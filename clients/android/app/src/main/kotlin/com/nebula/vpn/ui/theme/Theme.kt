package com.nebula.vpn.ui.theme

import android.os.Build
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

// Clean, Apple-inspired color palette
val NebulaPrimary = Color(0xFF000000)
val NebulaOnPrimary = Color(0xFFFFFFFF)
val NebulaPrimaryContainer = Color(0xFFF0F0F0)
val NebulaOnPrimaryContainer = Color(0xFF1C1C1E)

val NebulaSecondary = Color(0xFF8E8E93)
val NebulaOnSecondary = Color(0xFFFFFFFF)

val NebulaSurfaceLight = Color(0xFFF9F9F9)
val NebulaSurfaceDark = Color(0xFF000000)
val NebulaBackgroundLight = Color(0xFFFFFFFF)
val NebulaBackgroundDark = Color(0xFF0A0A0A)

val NebulaConnected = Color(0xFF34C759)
val NebulaDisconnected = Color(0xFF8E8E93)
val NebulaAccent = Color(0xFF007AFF)

private val LightColorScheme = lightColorScheme(
    primary = NebulaPrimary,
    onPrimary = NebulaOnPrimary,
    primaryContainer = NebulaPrimaryContainer,
    onPrimaryContainer = NebulaOnPrimaryContainer,
    secondary = NebulaSecondary,
    onSecondary = NebulaOnSecondary,
    surface = NebulaSurfaceLight,
    onSurface = Color(0xFF1C1C1E),
    background = NebulaBackgroundLight,
    onBackground = Color(0xFF1C1C1E),
    surfaceVariant = Color(0xFFF2F2F7),
    outline = Color(0xFFC6C6C8)
)

private val DarkColorScheme = darkColorScheme(
    primary = Color(0xFFFFFFFF),
    onPrimary = Color(0xFF000000),
    primaryContainer = Color(0xFF1C1C1E),
    onPrimaryContainer = Color(0xFFF0F0F0),
    secondary = Color(0xFF8E8E93),
    onSecondary = Color(0xFF000000),
    surface = Color(0xFF0A0A0A),
    onSurface = Color(0xFFF0F0F0),
    background = Color(0xFF000000),
    onBackground = Color(0xFFF0F0F0),
    surfaceVariant = Color(0xFF1C1C1E),
    outline = Color(0xFF38383A)
)

@Composable
fun NebulaVPNTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    dynamicColor: Boolean = false,
    content: @Composable () -> Unit
) {
    val colorScheme = when {
        dynamicColor && Build.VERSION.SDK_INT >= Build.VERSION_CODES.S -> {
            val context = androidx.compose.ui.platform.LocalContext.current
            if (darkTheme) dynamicDarkColorScheme(context)
            else dynamicLightColorScheme(context)
        }
        darkTheme -> DarkColorScheme
        else -> LightColorScheme
    }

    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography(
            displayLarge = Typography().displayLarge.copy(letterSpacing = androidx.compose.ui.unit.TextUnit(0f, androidx.compose.ui.unit.TextUnitType.Sp)),
            headlineLarge = Typography().headlineLarge.copy(letterSpacing = androidx.compose.ui.unit.TextUnit(0f, androidx.compose.ui.unit.TextUnitType.Sp)),
            titleLarge = Typography().titleLarge.copy(letterSpacing = androidx.compose.ui.unit.TextUnit(0f, androidx.compose.ui.unit.TextUnitType.Sp)),
            bodyLarge = Typography().bodyLarge.copy(letterSpacing = androidx.compose.ui.unit.TextUnit(0f, androidx.compose.ui.unit.TextUnitType.Sp)),
            labelLarge = Typography().labelLarge.copy(letterSpacing = androidx.compose.ui.unit.TextUnit(0.5f, androidx.compose.ui.unit.TextUnitType.Sp))
        ),
        content = content
    )
}
