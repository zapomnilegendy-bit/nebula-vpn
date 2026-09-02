package com.nebula.vpn

import android.app.Application
import com.nebula.vpn.data.remote.CrashReporter
import dagger.hilt.android.HiltAndroidApp

@HiltAndroidApp
class NebulaApp : Application() {
    override fun onCreate() {
        super.onCreate()
        CrashReporter.init(this)
    }
}
