package com.nebula.vpn.ui.screens

import androidx.compose.animation.*
import androidx.compose.animation.core.*
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.StrokeCap
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.nebula.vpn.ui.theme.NebulaConnected
import com.nebula.vpn.ui.theme.NebulaDisconnected

@Composable
fun MainScreen(
    selectedServer: VPNServer?,
    onServerPickerClick: () -> Unit,
    onSubscriptionClick: () -> Unit = {}
) {
    var isConnected by remember { mutableStateOf(false) }

    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(MaterialTheme.colorScheme.background)
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(horizontal = 24.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            // Top bar with subscription
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(top = 12.dp),
                horizontalArrangement = Arrangement.End
            ) {
                Surface(
                    onClick = onSubscriptionClick,
                    shape = RoundedCornerShape(8.dp),
                    color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.5f)
                ) {
                    Row(
                        modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp),
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        Text("★", fontSize = 12.sp)
                        Text(
                            "Pro",
                            fontSize = 12.sp,
                            fontWeight = FontWeight.Semibold
                        )
                    }
                }
            }

            Spacer().weight(0.5f)

            // Connection ring
            ConnectionRing(
                isConnected = isConnected,
                modifier = Modifier.size(220.dp)
            )

            Spacer().height(32.dp)

            // Status
            Text(
                text = if (isConnected) "Connected" else "Not Connected",
                fontSize = 20.sp,
                fontWeight = FontWeight.Semibold,
                color = MaterialTheme.colorScheme.onBackground
            )

            if (isConnected) {
                Text(
                    text = "0m 42s",
                    fontSize = 14.sp,
                    fontFamily = FontFamily.Monospace,
                    color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.5f),
                    modifier = Modifier.padding(top = 4.dp)
                )
            }

            Spacer().height(24.dp)

            // Data usage
            AnimatedVisibility(visible = isConnected) {
                Row(
                    horizontalArrangement = Arrangement.spacedBy(24.dp)
                ) {
                    DataBadge(icon = "↓", value = "1.2 MB")
                    DataBadge(icon = "↑", value = "342 KB")
                }
            }

            Spacer().weight(0.5f)

            // Server selector
            Surface(
                onClick = onServerPickerClick,
                shape = RoundedCornerShape(14.dp),
                color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.5f),
                modifier = Modifier.fillMaxWidth()
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp, vertical = 14.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = selectedServer?.country ?: "🌍",
                        fontSize = 22.sp
                    )
                    Spacer().width(12.dp)
                    Column(modifier = Modifier.weight(1f)) {
                        Text(
                            text = selectedServer?.city ?: "Select server",
                            fontSize = 15.sp,
                            fontWeight = FontWeight.Semibold,
                            color = MaterialTheme.colorScheme.onSurface
                        )
                        if (selectedServer != null) {
                            Text(
                                text = "${selectedServer!!.latencyMs}ms · ${selectedServer!!.loadPercent}% load",
                                fontSize = 12.sp,
                                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                            )
                        }
                    }
                    Text(
                        text = "›",
                        fontSize = 20.sp,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.3f)
                    )
                }
            }

            Spacer().height(12.dp)

            // Connect button
            Button(
                onClick = { isConnected = !isConnected },
                modifier = Modifier
                    .fillMaxWidth()
                    .height(50.dp),
                shape = RoundedCornerShape(14.dp),
                colors = ButtonDefaults.buttonColors(
                    containerColor = if (isConnected) Color(0xFFFF3B30)
                    else MaterialTheme.colorScheme.primary,
                    contentColor = if (isConnected) Color.White
                    else MaterialTheme.colorScheme.onPrimary
                )
            ) {
                Text(
                    text = if (isConnected) "Disconnect" else "Connect",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Semibold
                )
            }

            Spacer().height(24.dp)
        }
    }
}

@Composable
private fun ConnectionRing(isConnected: Boolean, modifier: Modifier = Modifier) {
    val infiniteTransition = rememberInfiniteTransition(label = "ring")
    val rotation by infiniteTransition.animateFloat(
        initialValue = 0f,
        targetValue = 360f,
        animationSpec = infiniteRepeatable(
            animation = tween(4000, easing = LinearEasing),
            repeatMode = RepeatMode.Restart
        ),
        label = "rotation"
    )

    val trimEnd by animateFloatAsState(
        targetValue = if (isConnected) 0.85f else 0.15f,
        animationSpec = spring(dampingRatio = 0.7f, stiffness = 100f),
        label = "trim"
    )

    Box(contentAlignment = Alignment.Center, modifier = modifier) {
        // Glow
        if (isConnected) {
            Box(
                modifier = Modifier
                    .size(260.dp)
                    .clip(CircleShape)
                    .background(NebulaConnected.copy(alpha = 0.08f))
            )
        }

        // Background ring
        Canvas(modifier = Modifier.fillMaxSize()) {
            drawCircle(
                color = Color.LightGray.copy(alpha = 0.15f),
                radius = size.minDimension / 2 * 0.92f,
                center = center,
                style = androidx.compose.ui.graphics.drawscope.Stroke(width = 3.dp.toPx())
            )
        }

        // Active ring
        Canvas(modifier = Modifier.fillMaxSize()) {
            val sweepAngle = trimEnd * 360f
            val arcSize = size.copy()

            drawArc(
                brush = Brush.sweepGradient(
                    colors = if (isConnected)
                        listOf(NebulaConnected, Color(0xFF007AFF))
                    else
                        listOf(NebulaDisconnected.copy(alpha = 0.2f), NebulaDisconnected.copy(alpha = 0.1f)),
                    center = Offset(size.width / 2, size.height / 2)
                ),
                startAngle = -90f + rotation,
                sweepAngle = sweepAngle,
                useCenter = false,
                style = androidx.compose.ui.graphics.drawscope.Stroke(
                    width = 3.dp.toPx(),
                    cap = StrokeCap.Round
                ),
                topLeft = Offset(arcSize.width * 0.05f, arcSize.height * 0.05f),
                size = androidx.compose.ui.geometry.Size(arcSize.width * 0.9f, arcSize.height * 0.9f)
            )
        }

        // Center icon
        Text(
            text = if (isConnected) "🛡️" else "🔒",
            fontSize = 48.sp
        )
    }
}

@Composable
private fun DataBadge(icon: String, value: String) {
    Surface(
        shape = RoundedCornerShape(10.dp),
        color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.5f)
    ) {
        Row(
            modifier = Modifier.padding(horizontal = 14.dp, vertical = 8.dp),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(6.dp)
        ) {
            Text(
                text = icon,
                fontSize = 12.sp,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
            )
            Text(
                text = value,
                fontSize = 13.sp,
                fontWeight = FontWeight.Medium,
                fontFamily = FontFamily.Monospace,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
            )
        }
    }
}
