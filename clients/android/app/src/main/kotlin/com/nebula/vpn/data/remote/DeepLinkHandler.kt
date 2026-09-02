package com.nebula.vpn.data.remote

import android.content.Context
import android.content.Intent
import android.net.Uri
import android.util.Log

enum class DeepLinkRoute {
    CONNECT, SETTINGS, SUBSCRIPTION, PROMO, UNKNOWN
}

data class DeepLinkResult(
    val route: DeepLinkRoute,
    val serverId: String? = null,
    val promoCode: String? = null
)

object DeepLinkHandler {
    private const val TAG = "DeepLinkHandler"
    private const val SCHEME = "nebula"
    private const val UNIVERSAL_HOST = "nebula-vpn.io"

    fun handle(intent: Intent): DeepLinkResult? {
        val uri = intent.data ?: return null
        return handleUri(uri)
    }

    fun handleUri(uri: Uri): DeepLinkResult? {
        Log.d(TAG, "Handling URI: $uri")

        return when (uri.scheme?.lowercase()) {
            SCHEME -> handleNebulaScheme(uri)
            "https", "http" -> handleUniversalLink(uri)
            else -> {
                Log.w(TAG, "Unknown scheme: ${uri.scheme}")
                null
            }
        }
    }

    private fun handleNebulaScheme(uri: Uri): DeepLinkResult {
        return when (uri.host?.lowercase()) {
            "connect" -> DeepLinkResult(DeepLinkRoute.CONNECT)
            "server" -> DeepLinkResult(
                DeepLinkRoute.CONNECT,
                serverId = uri.getQueryParameter("id")
            )
            "settings" -> DeepLinkResult(DeepLinkRoute.SETTINGS)
            "subscription" -> DeepLinkResult(DeepLinkRoute.SUBSCRIPTION)
            "promo" -> DeepLinkResult(
                DeepLinkRoute.PROMO,
                promoCode = uri.getQueryParameter("code")
            )
            else -> DeepLinkResult(DeepLinkRoute.UNKNOWN)
        }
    }

    private fun handleUniversalLink(uri: Uri): DeepLinkResult {
        val path = uri.path?.lowercase() ?: ""
        return when {
            path == "/connect" -> DeepLinkResult(DeepLinkRoute.CONNECT)
            path.startsWith("/server/") -> DeepLinkResult(
                DeepLinkRoute.CONNECT,
                serverId = path.removePrefix("/server/")
            )
            path == "/settings" -> DeepLinkResult(DeepLinkRoute.SETTINGS)
            path == "/subscription" -> DeepLinkResult(DeepLinkRoute.SUBSCRIPTION)
            path.startsWith("/promo/") -> DeepLinkResult(
                DeepLinkRoute.PROMO,
                promoCode = path.removePrefix("/promo/")
            )
            else -> DeepLinkResult(DeepLinkRoute.UNKNOWN)
        }
    }

    fun createDeepLinkIntent(context: Context, route: DeepLinkRoute, serverId: String? = null): Intent {
        val uri = Uri.Builder()
            .scheme(SCHEME)
            .authority(
                when (route) {
                    DeepLinkRoute.CONNECT -> "connect"
                    DeepLinkRoute.SETTINGS -> "settings"
                    DeepLinkRoute.SUBSCRIPTION -> "subscription"
                    DeepLinkRoute.PROMO -> "promo"
                    DeepLinkRoute.UNKNOWN -> ""
                }
            )
            .build()

        return Intent(Intent.ACTION_VIEW, uri, context, Class.forName("com.nebula.vpn.MainActivity"))
    }
}
