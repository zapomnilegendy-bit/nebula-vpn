package com.nebula.vpn.vpn

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Intent
import android.net.VpnService
import android.os.Build
import android.os.ParcelFileDescriptor
import com.nebula.vpn.MainActivity

class NebulaVpnService : VpnService() {

    private var tunnelInterface: ParcelFileDescriptor? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        startForeground(NOTIFICATION_ID, createNotification())

        if (intent?.action == ACTION_DISCONNECT) {
            disconnect()
            return START_NOT_STICKY
        }

        connect()
        return START_STICKY
    }

    private fun connect() {
        val builder = Builder()
            .setSession("NebulaVPN")
            .setMtu(1420)
            .addAddress("10.0.0.2", 32)
            .addDnsServer("1.1.1.1")
            .addDnsServer("1.0.0.1")
            .addRoute("0.0.0.0", 0)

        tunnelInterface?.close()
        tunnelInterface = builder.establish()
    }

    private fun disconnect() {
        tunnelInterface?.close()
        tunnelInterface = null
        stopForeground(STOP_FOREGROUND_REMOVE)
        stopSelf()
    }

    private fun createNotification(): Notification {
        val channelId = "nebula_vpn_channel"

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                channelId,
                "NebulaVPN",
                NotificationManager.IMPORTANCE_LOW
            )
            val manager = getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(channel)
        }

        val disconnectIntent = Intent(this, NebulaVpnService::class.java).apply {
            action = ACTION_DISCONNECT
        }

        return Notification.Builder(this, channelId)
            .setContentTitle("NebulaVPN")
            .setContentText("Connected - Securing your traffic")
            .setSmallIcon(android.R.drawable.ic_lock_lock)
            .setOngoing(true)
            .addAction(
                android.R.drawable.ic_media_pause,
                "Disconnect",
                PendingIntent.getService(this, 0, disconnectIntent, PendingIntent.FLAG_IMMUTABLE)
            )
            .build()
    }

    override fun onDestroy() {
        disconnect()
        super.onDestroy()
    }

    companion object {
        const val ACTION_DISCONNECT = "com.nebula.vpn.DISCONNECT"
        private const val NOTIFICATION_ID = 1
    }
}
