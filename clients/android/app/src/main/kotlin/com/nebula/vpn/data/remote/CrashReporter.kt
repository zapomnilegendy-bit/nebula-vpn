package com.nebula.vpn.data.remote

import android.content.Context
import android.os.Build
import android.util.Log
import java.io.PrintWriter
import java.io.StringWriter
import java.util.UUID

object CrashReporter {
    private const val TAG = "CrashReporter"
    private const val API_URL = "https://api.nebula-vpn.io/v1/crashes"
    private var enabled = false
    private var deviceId: String = UUID.randomUUID().toString().take(8)

    fun init(context: Context) {
        deviceId = "android-${deviceId}"
        enabled = true
        setupExceptionHandler()
        Log.d(TAG, "CrashReporter initialized on ${Build.MODEL} (API ${Build.VERSION.SDK_INT})")
    }

    private fun setupExceptionHandler() {
        val defaultHandler = Thread.getDefaultUncaughtExceptionHandler()

        Thread.setDefaultUncaughtExceptionHandler { thread, throwable ->
            val report = buildReport(throwable)
            sendCrashSync(report)
            defaultHandler?.uncaughtException(thread, throwable)
        }
    }

    fun report(throwable: Throwable, context: Map<String, Any> = emptyMap()) {
        if (!enabled) return
        val report = buildReport(throwable, context)
        sendCrashAsync(report)
    }

    fun reportMessage(message: String, level: String = "error", context: Map<String, Any> = emptyMap()) {
        if (!enabled) return
        val payload = mutableMapOf<String, Any>(
            "message" to message,
            "level" to level,
            "platform" to "android",
            "device_id" to deviceId
        )
        payload.putAll(context)
        sendCrashAsync(payload)
    }

    private fun buildReport(throwable: Throwable, extra: Map<String, Any> = emptyMap()): Map<String, Any> {
        val sw = StringWriter()
        throwable.printStackTrace(PrintWriter(sw))

        return mutableMapOf(
            "exception" to throwable.javaClass.name,
            "message" to (throwable.message ?: ""),
            "stacktrace" to sw.toString(),
            "platform" to "android",
            "device_id" to deviceId,
            "device_model" to Build.MODEL,
            "device_brand" to Build.BRAND,
            "os_version" to Build.VERSION.RELEASE,
            "api_level" to Build.VERSION.SDK_INT,
            "version" to "1.0.0",
        ).apply { putAll(extra) }
    }

    private fun sendCrashSync(payload: Map<String, Any>) {
        try {
            val json = org.json.JSONObject(payload).toString()
            val url = java.net.URL(API_URL)
            val conn = url.openConnection() as java.net.HttpURLConnection
            conn.requestMethod = "POST"
            conn.setRequestProperty("Content-Type", "application/json")
            conn.doOutput = true
            conn.outputStream.write(json.toByteArray())
            conn.responseCode
            conn.disconnect()
        } catch (_: Exception) {}
    }

    private fun sendCrashAsync(payload: Map<String, Any>) {
        Thread { sendCrashSync(payload) }.start()
    }
}
