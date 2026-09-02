package com.nebula.vpn

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import com.nebula.vpn.data.remote.CrashReporter
import com.nebula.vpn.data.remote.DeepLinkHandler
import com.nebula.vpn.data.remote.DeepLinkRoute
import com.nebula.vpn.ui.screens.MainScreen
import com.nebula.vpn.ui.screens.OnboardingScreen
import com.nebula.vpn.ui.screens.ServerListScreen
import com.nebula.vpn.ui.screens.SubscriptionScreen
import com.nebula.vpn.ui.screens.VPNServer
import com.nebula.vpn.ui.theme.NebulaVPNTheme
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    private val prefs by lazy { getSharedPreferences("nebula", MODE_PRIVATE) }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val hasSeenOnboarding = prefs.getBoolean("onboarding_done", false)

        setContent {
            NebulaVPNTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    NebulaAppContent(
                        hasSeenOnboarding = hasSeenOnboarding,
                        onOnboardingComplete = {
                            prefs.edit().putBoolean("onboarding_done", true).apply()
                        },
                        deepLinkRoute = parseDeepLink(intent)
                    )
                }
            }
        }
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        // Handle deep links when app is already open
        val route = DeepLinkHandler.handle(intent)
        if (route != null) {
            handleDeepLinkRoute(route.route, route.serverId)
        }
    }

    private fun parseDeepLink(intent: Intent): Pair<DeepLinkRoute, String?>? {
        val result = DeepLinkHandler.handle(intent) ?: return null
        return result.route to result.serverId
    }

    private fun handleDeepLinkRoute(route: DeepLinkRoute, serverId: String?) {
        when (route) {
            DeepLinkRoute.SETTINGS -> {
                CrashReporter.reportMessage("Deep link: settings", level = "info")
            }
            DeepLinkRoute.CONNECT -> {
                CrashReporter.reportMessage("Deep link: connect", level = "info")
            }
            DeepLinkRoute.SUBSCRIPTION -> {
                CrashReporter.reportMessage("Deep link: subscription", level = "info")
            }
            DeepLinkRoute.PROMO -> {
                CrashReporter.reportMessage("Deep link: promo", level = "info")
            }
            DeepLinkRoute.UNKNOWN -> {}
        }
    }
}

@Composable
private fun NebulaAppContent(
    hasSeenOnboarding: Boolean,
    onOnboardingComplete: () -> Unit,
    deepLinkRoute: Pair<DeepLinkRoute, String?>?
) {
    var showOnboarding by remember { mutableStateOf(!hasSeenOnboarding) }
    var showServerList by remember { mutableStateOf(false) }
    var showSubscription by remember { mutableStateOf(
        deepLinkRoute?.first == DeepLinkRoute.SUBSCRIPTION
    ) }
    var selectedServer by remember { mutableStateOf<VPNServer?>(null) }
    var deepLinkHandled by remember { mutableStateOf(false) }

    // Handle deep link route
    LaunchedEffect(deepLinkRoute) {
        if (!deepLinkHandled && deepLinkRoute != null) {
            when (deepLinkRoute.first) {
                DeepLinkRoute.CONNECT -> showServerList = true
                DeepLinkRoute.SETTINGS -> { /* Open settings */ }
                DeepLinkRoute.SUBSCRIPTION -> showSubscription = true
                else -> {}
            }
            deepLinkHandled = true
        }
    }

    when {
        showOnboarding -> {
            OnboardingScreen(
                onComplete = {
                    showOnboarding = false
                    onOnboardingComplete()
                }
            )
        }
        showSubscription -> {
            SubscriptionScreen(onDismiss = { showSubscription = false })
        }
        showServerList -> {
            ServerListScreen(
                selectedServerId = selectedServer?.id ?: deepLinkRoute?.second,
                onServerSelected = { server ->
                    selectedServer = server
                    showServerList = false
                },
                onDismiss = { showServerList = false }
            )
        }
        else -> {
            MainScreen(
                selectedServer = selectedServer,
                onServerPickerClick = { showServerList = true },
                onSubscriptionClick = { showSubscription = true }
            )
        }
    }
}
