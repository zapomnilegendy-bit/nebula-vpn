package com.nebula.vpn.ui.screens

import androidx.compose.animation.*
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.pager.HorizontalPager
import androidx.compose.foundation.pager.rememberPagerState
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import kotlinx.coroutines.launch

data class OnboardingPage(
    val icon: String,
    val title: String,
    val description: String
)

val onboardingPages = listOf(
    OnboardingPage("🛡️", "Privacy First",
        "NebulaVPN encrypts all your traffic with military-grade WireGuard protocol. No logs, no tracking."),
    OnboardingPage("🌍", "Global Network",
        "Access content from 50+ locations worldwide. Your real IP stays hidden everywhere you go."),
    OnboardingPage("⚡", "Anti-DPI Technology",
        "Advanced obfuscation bypasses deep packet inspection. Works in the most restrictive networks."),
    OnboardingPage("📱", "One-Click Connect",
        "Import config via QR code or select from our optimized server list. Connect in one click."),
)

@Composable
fun OnboardingScreen(onComplete: () -> Unit) {
    val pagerState = rememberPagerState(pageCount = { onboardingPages.size })
    val scope = rememberCoroutineScope()

    Column(
        modifier = Modifier.fillMaxSize(),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        HorizontalPager(
            state = pagerState,
            modifier = Modifier.weight(1f)
        ) { page ->
            val p = onboardingPages[page]
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(horizontal = 40.dp),
                verticalArrangement = Arrangement.Center,
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = p.icon,
                    fontSize = 72.sp
                )
                Spacer().height(32.dp)
                Text(
                    text = p.title,
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = MaterialTheme.colorScheme.onBackground
                )
                Spacer().height(12.dp)
                Text(
                    text = p.description,
                    fontSize = 15.sp,
                    color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.6f),
                    textAlign = TextAlign.Center,
                    lineHeight = 22.sp
                )
            }
        }

        // Page indicator
        Row(
            horizontalArrangement = Arrangement.spacedBy(8.dp),
            modifier = Modifier.padding(bottom = 16.dp)
        ) {
            repeat(onboardingPages.size) { i ->
                Box(
                    modifier = Modifier
                        .size(if (i == pagerState.currentPage) 24.dp else 8.dp, 8.dp)
                        .then(
                            Modifier
                                .let {
                                    if (i == pagerState.currentPage)
                                        androidx.compose.ui.draw.clip(androidx.compose.foundation.shape.RoundedCornerShape(4.dp))
                                    else
                                        androidx.compose.ui.draw.clip(androidx.compose.foundation.shape.CircleShape)
                                }
                        )
                        .then(
                            Modifier.let {
                                androidx.compose.material3.MaterialTheme.colorScheme.primary.let { color ->
                                    Modifier.background(
                                        if (i == pagerState.currentPage) color
                                        else color.copy(alpha = 0.2f)
                                    )
                                }
                            }
                        )
                )
            }
        }

        // Buttons
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 24.dp, vertical = 16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            TextButton(onClick = onComplete) {
                Text("Skip", color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.5f))
            }

            Button(
                onClick = {
                    if (pagerState.currentPage < onboardingPages.size - 1) {
                        scope.launch {
                            pagerState.animateScrollToPage(pagerState.currentPage + 1)
                        }
                    } else {
                        onComplete()
                    }
                }
            ) {
                Text(
                    if (pagerState.currentPage < onboardingPages.size - 1) "Continue"
                    else "Get Started"
                )
            }
        }

        Spacer().height(8.dp)
    }
}
