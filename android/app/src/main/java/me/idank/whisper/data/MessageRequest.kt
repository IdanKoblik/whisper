package me.idank.whisper.data

import com.fasterxml.jackson.annotation.JsonProperty

data class MessageRequest(
    @field:JsonProperty("device_id") val deviceId: String,
    val message: String,
    val subscribers: List<String>
)
