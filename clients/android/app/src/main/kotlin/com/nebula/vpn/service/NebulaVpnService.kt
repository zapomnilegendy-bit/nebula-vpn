package com.nebula.vpn.service

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Intent
import android.net.VpnService
import android.os.Build
import android.os.ParcelFileDescriptor
import kotlinx.coroutines.*
import java.io.FileInputStream
import java.io.FileOutputStream
import java.net.InetSocketAddress
import java.nio.ByteBuffer

class NebulaVpnService : VpnService() {

    private var vpnInterface: ParcelFileDescriptor? = null
    private var isRunning = false
    private val serviceScope = CoroutineScope(Dispatchers.IO + SupervisorJob())
    private val inputBuffer = ByteBuffer.allocateDirect(32767)
    private val outputBuffer = ByteBuffer.allocateDirect(32767)
    private var idleCount = 0

    companion object {
        const val ACTION_CONNECT = "com.nebula.vpn.CONNECT"
        const val ACTION_DISCONNECT = "com.nebula.vpn.DISCONNECT"
        const val CHANNEL_ID = "nebula_vpn_channel"
        const val NOTIFICATION_ID = 1
    }

    override fun onCreate() {
        super.onCreate()
        createNotificationChannel()
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        when (intent?.action) {
            ACTION_CONNECT -> {
                val server = intent.getStringExtra("server_host") ?: return START_STICKY
                val port = intent.getIntExtra("server_port", 51820)
                val config = intent.getStringExtra("config") ?: return START_STICKY
                startVpn(server, port, config)
            }
            ACTION_DISCONNECT -> stopVpn()
        }
        return START_STICKY
    }

    override fun onDestroy() {
        super.onDestroy()
        stopVpn()
        serviceScope.cancel()
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                "Nebula VPN",
                NotificationManager.IMPORTANCE_LOW
            ).apply {
                description = "VPN connection status"
                setShowBadge(false)
            }
            val manager = getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(channel)
        }
    }

    private fun startVpn(serverHost: String, serverPort: Int, config: String) {
        if (isRunning) stopVpn()

        val builder = Builder()

        // Battery-efficient tunnel setup
        builder.setSession("Nebula VPN")
        builder.setMtu(1400)

        // Add virtual addresses
        builder.addAddress("10.0.0.2", 32)

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
            builder.setMetered(false)
            builder.setUnderlyingNetworks(null) // Let system manage networks
        }

        // Route all traffic
        builder.addRoute("0.0.0.0", 0)
        builder.addRoute("::", 0)

        // DNS
        builder.addDnsServer("1.1.1.1")
        builder.addDnsServer("8.8.8.8")
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
            builder.setBlocking(true)
        }

        // Add allowed applications (split tunnel)
        builder.addAllowedApplication(packageName)

        // Establish VPN
        vpnInterface = builder.establish()

        if (vpnInterface == null) {
            stopForeground(STOP_FOREGROUND_DETACH)
            return
        }

        // Start foreground notification
        startForeground(NOTIFICATION_ID, createNotification())

        isRunning = true

        // Start packet processing
        serviceScope.launch {
            processPackets()
        }
    }

    private fun stopVpn() {
        isRunning = false
        try {
            vpnInterface?.close()
        } catch (_: Exception) {}
        vpnInterface = null
        stopForeground(STOP_FOREGROUND_REMOVE)
        stopSelf()
    }

    private suspend fun processPackets() {
        val input = FileInputStream(vpnInterface?.fileDescriptor)
        val output = FileOutputStream(vpnInterface?.fileDescriptor)

        while (isRunning) {
            try {
                // Read with timeout to avoid busy-waiting
                val read = withTimeout(1000L) {
                    input.channel.read(inputBuffer)
                }

                if (read < 0) break

                // Battery optimization: idle detection
                if (read == 0) {
                    idleCount++
                    if (idleCount > 50) {
                        delay(200) // Sleep 200ms when idle
                    }
                    continue
                }
                idleCount = 0

                inputBuffer.flip()
                val packet = ByteArray(read)
                inputBuffer.get(packet)
                inputBuffer.clear()

                // Encrypt and send
                val encrypted = processOutgoing(packet)
                outputBuffer.put(encrypted)
                outputBuffer.flip()
                output.channel.write(outputBuffer)
                outputBuffer.clear()

                // Yield to avoid hogging CPU
                if (idleCount % 10 == 0) yield()

            } catch (_: CancellationException) {
                break
            } catch (e: Exception) {
                if (isRunning) {
                    delay(1000)
                }
            }
        }
    }

    private fun processOutgoing(packet: ByteArray): ByteArray {
        // In production: encrypt via Rust core
        // Add traffic padding to defeat analysis
        val padding = if (packet.size < 200) {
            ByteArray((40..200).random()) { 0 }
        } else {
            ByteArray(0)
        }
        return packet + padding
    }

    private fun createNotification(): Notification {
        val builder = Notification.Builder(this, CHANNEL_ID)
            .setContentTitle("Nebula VPN")
            .setContentText("Connected • Securing your traffic")
            .setSmallIcon(android.R.drawable.ic_lock_lock)
            .setOngoing(true)
            .setPriority(Notification.PRIORITY_LOW)

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            builder.setChannelId(CHANNEL_ID)
        }

        return builder.build()
    }

    override fun onRevoke() {
        stopVpn()
    }
}
