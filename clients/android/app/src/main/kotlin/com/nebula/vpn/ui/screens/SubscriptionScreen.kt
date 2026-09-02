package com.nebula.vpn.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

data class SubscriptionPlan(
    val id: String,
    val name: String,
    val price: String,
    val period: String,
    val features: List<String>,
    val isPopular: Boolean = false,
    val isCurrent: Boolean = false
)

val subscriptionPlans = listOf(
    SubscriptionPlan(
        id = "free", name = "Free", price = "$0", period = "forever",
        features = listOf(
            "3 server locations",
            "Basic speed (10 Mbps)",
            "1 device",
            "Community support"
        )
    ),
    SubscriptionPlan(
        id = "pro", name = "Pro", price = "$9.99", period = "per month",
        features = listOf(
            "50+ server locations",
            "Unlimited speed",
            "5 devices",
            "Anti-DPI technology",
            "Priority support",
            "No logs"
        ),
        isPopular = true,
        isCurrent = true
    ),
    SubscriptionPlan(
        id = "ultimate", name = "Ultimate", price = "$19.99", period = "per month",
        features = listOf(
            "Everything in Pro",
            "Unlimited devices",
            "Dedicated IP",
            "Multi-hop VPN",
            "24/7 priority support",
            "Family sharing (5 accounts)"
        )
    )
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SubscriptionScreen(onDismiss: () -> Unit) {
    var selectedPlan by remember { mutableStateOf("pro") }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        "Subscription",
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
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(horizontal = 20.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Spacer().height(8.dp)

            // Current plan badge
            Card(
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.primaryContainer
                ),
                modifier = Modifier.fillMaxWidth()
            ) {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(20.dp),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Text(
                        "Pro Plan",
                        fontSize = 20.sp,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.onPrimaryContainer
                    )
                    Spacer().height(4.dp)
                    Text(
                        "Your subscription is active",
                        fontSize = 14.sp,
                        color = MaterialTheme.colorScheme.onPrimaryContainer.copy(alpha = 0.7f)
                    )
                    Spacer().height(12.dp)
                    Text(
                        "$9.99",
                        fontSize = 36.sp,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.onPrimaryContainer
                    )
                    Text(
                        "per month · renews Jun 20, 2026",
                        fontSize = 12.sp,
                        color = MaterialTheme.colorScheme.onPrimaryContainer.copy(alpha = 0.6f)
                    )
                }
            }

            Spacer().height(24.dp)

            // Plan options
            Text(
                "Choose your plan",
                fontSize = 13.sp,
                fontWeight = FontWeight.Semibold,
                letterSpacing = 0.5.sp,
                color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.5f),
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(bottom = 12.dp)
            )

            subscriptionPlans.forEach { plan ->
                PlanCard(
                    plan = plan,
                    isSelected = plan.id == selectedPlan,
                    onClick = { selectedPlan = plan.id }
                )
                Spacer().height(8.dp)
            }

            Spacer().weight(1f)

            // Subscribe button
            Button(
                onClick = { /* Start Stripe/Google IAP flow */ },
                modifier = Modifier
                    .fillMaxWidth()
                    .height(50.dp),
                shape = RoundedCornerShape(14.dp)
            ) {
                Text(
                    "Continue with ${selectedPlan.replaceFirstChar { it.uppercase() }}",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Semibold
                )
            }

            Spacer().height(8.dp)

            Text(
                "Cancel anytime. Auto-renews unless cancelled.",
                fontSize = 11.sp,
                color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.4f),
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(bottom = 24.dp)
            )
        }
    }
}

@Composable
private fun PlanCard(
    plan: SubscriptionPlan,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    Card(
        onClick = onClick,
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(
            containerColor = if (isSelected)
                MaterialTheme.colorScheme.primary.copy(alpha = 0.06f)
            else MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.3f)
        ),
        border = if (isSelected)
            CardDefaults.outlinedCardBorder().copy(
                brush = androidx.compose.ui.graphics.SolidColor(
                    MaterialTheme.colorScheme.primary
                )
            )
        else null,
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                modifier = Modifier.fillMaxWidth()
            ) {
                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        plan.name,
                        fontSize = 16.sp,
                        fontWeight = FontWeight.Semibold
                    )
                    Row(verticalAlignment = Alignment.Bottom) {
                        Text(
                            plan.price,
                            fontSize = 24.sp,
                            fontWeight = FontWeight.Bold
                        )
                        Spacer().width(4.dp)
                        Text(
                            plan.period,
                            fontSize = 12.sp,
                            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                        )
                    }
                }

                if (plan.isPopular) {
                    Surface(
                        shape = RoundedCornerShape(6.dp),
                        color = MaterialTheme.colorScheme.primary
                    ) {
                        Text(
                            "POPULAR",
                            fontSize = 9.sp,
                            fontWeight = FontWeight.Bold,
                            color = MaterialTheme.colorScheme.onPrimary,
                            modifier = Modifier.padding(horizontal = 8.dp, vertical = 3.dp)
                        )
                    }
                }
            }

            Spacer().height(8.dp)

            plan.features.forEach { feature ->
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    modifier = Modifier.padding(vertical = 2.dp)
                ) {
                    Text(
                        "✓",
                        fontSize = 13.sp,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.primary
                    )
                    Spacer().width(8.dp)
                    Text(
                        feature,
                        fontSize = 13.sp,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
                    )
                }
            }
        }
    }
}
