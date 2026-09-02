# Nebula VPN ProGuard Rules

# Keep Rust JNI bridge
-keep class com.nebula.vpn.NebulaCore { *; }

# Keep Retrofit interfaces
-keep,allowobfuscation interface com.nebula.vpn.api.** { *; }
-keep class com.nebula.vpn.api.** { *; }

# Keep Hilt/DI
-keep class dagger.hilt.** { *; }
-keep class javax.inject.** { *; }
-keep class * extends dagger.hilt.android.internal.managers.ViewComponentManager$FragmentContextWrapper { *; }

# Keep WorkManager
-keep class androidx.work.** { *; }

# Keep data classes used by Gson
-keep class com.nebula.vpn.model.** { *; }
-keep class com.nebula.vpn.api.** { *; }

# Keep EncryptedSharedPreferences
-keep class androidx.security.crypto.** { *; }

# Keep notification builder
-keep class com.nebula.vpn.service.NebulaVpnService { *; }

# AndroidX
-dontwarn androidx.**

# OkHttp
-dontwarn okhttp3.**
-dontwarn okio.**

# Retrofit
-dontwarn retrofit2.**
-keep class retrofit2.** { *; }

# Coroutines
-dontwarn kotlinx.coroutines.**
-keepnames class kotlinx.coroutines.** { *; }

# Gson
-keepclassmembers class * {
    @com.google.gson.annotations.SerializedName <fields>;
}

# Keep ViewModels
-keep class * extends androidx.lifecycle.ViewModel { *; }

# Generic rules
-keepattributes Signature
-keepattributes *Annotation*
-keepattributes EnclosingMethod
-keepattributes InnerClasses

# Remove logging in release
-assumenosideeffects class android.util.Log {
    public static boolean isLoggable(java.lang.String, int);
    public static int v(...);
    public static int d(...);
    public static int i(...);
}
