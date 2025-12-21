package dev.idank.whisper.requests;

import com.fasterxml.jackson.annotation.JsonProperty;

public record HeartbeatRequest(
        @JsonProperty("device")
        String deviceId
) {
}
