package com.nebula.vpn.data.api

import retrofit2.Response
import retrofit2.http.*

data class AuthRequest(
    val email: String,
    val password: String
)

data class AuthResponse(
    val access_token: String,
    val refresh_token: String,
    val expires_in: Long,
    val user_id: String
)

data class ServerResponse(
    val id: String,
    val name: String,
    val country: String,
    val city: String,
    val protocols: List<String>,
    val load_percent: Int,
    val latency_ms: Int,
    val is_healthy: Boolean
)

interface VpnApi {
    @POST("auth/login")
    suspend fun login(@Body request: AuthRequest): Response<AuthResponse>

    @POST("auth/register")
    suspend fun register(@Body request: AuthRequest): Response<AuthResponse>

    @GET("servers")
    suspend fun getServers(
        @Query("country") country: String? = null,
        @Query("protocol") protocol: String? = null
    ): Response<List<ServerResponse>>

    @GET("servers/recommend")
    suspend fun getRecommendedServer(): Response<ServerResponse>
}
