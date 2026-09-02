package com.nebula.vpn.ui.screens

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

data class VPNServer(
    val id: String,
    val name: String,
    val city: String,
    val country: String,
    val host: String,
    val port: Int,
    val latencyMs: Int,
    val loadPercent: Int,
    val isPremium: Boolean
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ServerListScreen(
    selectedServerId: String?,
    onServerSelected: (VPNServer) -> Unit,
    onDismiss: () -> Unit
) {
    val servers = remember { sampleServersAndroid }
    var searchQuery by remember { mutableStateOf("") }

    val filteredServers = remember(servers, searchQuery) {
        if (searchQuery.isBlank()) servers
        else servers.filter {
            it.city.contains(searchQuery, ignoreCase = true) ||
            it.name.contains(searchQuery, ignoreCase = true) ||
            it.country.contains(searchQuery, ignoreCase = true)
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        "Servers",
                        fontWeight = FontWeight.Semibold
                    )
                },
                actions = {
                    TextButton(onClick = onDismiss) {
                        Text("Done", fontWeight = FontWeight.Semibold)
                    }
                }
            )
        }
    ) { padding ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding),
            contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp),
            verticalArrangement = Arrangement.spacedBy(4.dp)
        ) {
            // Fastest server
            item {
                Surface(
                    onClick = {
                        val fastest = servers.filter { it.latencyMs > 0 }
                            .minByOrNull { it.latencyMs }
                        fastest?.let { onServerSelected(it) }
                    },
                    shape = RoundedCornerShape(12.dp),
                    color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.3f),
                    modifier = Modifier.fillMaxWidth()
                ) {
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(12.dp),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Surface(
                            shape = RoundedCornerShape(10.dp),
                            color = MaterialTheme.colorScheme.primary.copy(alpha = 0.1f),
                            modifier = Modifier.size(40.dp)
                        ) {
                            Box(contentAlignment = Alignment.Center) {
                                Text(
                                    "⚡",
                                    fontSize = 18.sp
                                )
                            }
                        }
                        Spacer().width(12.dp)
                        Column(modifier = Modifier.weight(1f)) {
                            Text(
                                "Fastest Server",
                                fontSize = 15.sp,
                                fontWeight = FontWeight.Semibold
                            )
                            Text(
                                "Automatically selected",
                                fontSize = 12.sp,
                                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                            )
                        }
                        if (selectedServerId == "auto") {
                            Text(
                                "✓",
                                fontSize = 16.sp,
                                fontWeight = FontWeight.Bold,
                                color = MaterialTheme.colorScheme.primary
                            )
                        }
                    }
                }
            }

            item { Spacer().height(8.dp) }

            // Section header
            item {
                Text(
                    "ALL LOCATIONS",
                    fontSize = 11.sp,
                    fontWeight = FontWeight.Semibold,
                    letterSpacing = 0.8.sp,
                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f),
                    modifier = Modifier.padding(vertical = 8.dp, horizontal = 4.dp)
                )
            }

            items(filteredServers) { server ->
                ServerCard(
                    server = server,
                    isSelected = server.id == selectedServerId,
                    onClick = { onServerSelected(server) }
                )
            }
        }
    }
}

@Composable
private fun ServerCard(
    server: VPNServer,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    Surface(
        onClick = onClick,
        shape = RoundedCornerShape(12.dp),
        color = if (isSelected)
            MaterialTheme.colorScheme.primary.copy(alpha = 0.06f)
        else Color.Transparent,
        modifier = Modifier.fillMaxWidth()
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                server.country,
                fontSize = 22.sp
            )
            Spacer().width(12.dp)
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    server.city,
                    fontSize = 14.sp,
                    fontWeight = FontWeight.Semibold
                )
                Row(
                    horizontalArrangement = Arrangement.spacedBy(4.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        server.name,
                        fontSize = 11.sp,
                        fontFamily = FontFamily.Monospace,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                    )
                    if (server.isPremium) {
                        Surface(
                            shape = RoundedCornerShape(3.dp),
                            color = Color(0xFFFF9500)
                        ) {
                            Text(
                                "PREMIUM",
                                fontSize = 8.sp,
                                fontWeight = FontWeight.Bold,
                                color = Color.White,
                                modifier = Modifier.padding(horizontal = 4.dp, vertical = 1.dp)
                            )
                        }
                    }
                }
            }
            Column(horizontalAlignment = Alignment.End) {
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    Box(
                        modifier = Modifier
                            .size(6.dp)
                            .clip(CircleShape)
                            .then(
                                Modifier.clip(CircleShape)
                            )
                            .let { mod ->
                                mod.then(
                                    Modifier.clip(CircleShape)
                                )
                            }
                    )
                    Text(
                        "${server.latencyMs}ms",
                        fontSize = 12.sp,
                        fontWeight = FontWeight.Medium,
                        fontFamily = FontFamily.Monospace,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                    )
                }
            }
            if (isSelected) {
                Spacer().width(8.dp)
                Text(
                    "✓",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Bold,
                    color = MaterialTheme.colorScheme.primary
                )
            }
        }
    }
}

val sampleServersAndroid = listOf(
    VPNServer("eu-fra-01", "fra1.nebula", "Frankfurt", "🇩🇪", "fra1.nebula-vpn.io", 51820, 12, 42, false),
    VPNServer("eu-ams-01", "ams1.nebula", "Amsterdam", "🇳🇱", "ams1.nebula-vpn.io", 51820, 18, 58, false),
    VPNServer("gb-lon-01", "lon1.nebula", "London", "🇬🇧", "lon1.nebula-vpn.io", 51820, 24, 55, false),
    VPNServer("us-nyc-01", "nyc1.nebula", "New York", "🇺🇸", "nyc1.nebula-vpn.io", 51820, 89, 71, true),
    VPNServer("us-sfo-01", "sfo1.nebula", "San Francisco", "🇺🇸", "sfo1.nebula-vpn.io", 51820, 105, 34, false),
    VPNServer("us-mia-01", "mia1.nebula", "Miami", "🇺🇸", "mia1.nebula-vpn.io", 51820, 130, 29, false),
    VPNServer("jp-tokyo-01", "tky1.nebula", "Tokyo", "🇯🇵", "tky1.nebula-vpn.io", 51820, 185, 23, true),
    VPNServer("sg-sin-01", "sin1.nebula", "Singapore", "🇸🇬", "sin1.nebula-vpn.io", 51820, 210, 18, false),
    VPNServer("au-syd-01", "syd1.nebula", "Sydney", "🇦🇺", "syd1.nebula-vpn.io", 51820, 280, 12, true),
    VPNServer("br-gru-01", "gru1.nebula", "São Paulo", "🇧🇷", "gru1.nebula-vpn.io", 51820, 175, 9, false),
)
