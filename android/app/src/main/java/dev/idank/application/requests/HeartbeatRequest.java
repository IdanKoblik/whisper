package dev.idank.application.requests;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public record HeartbeatRequest(
        @JsonProperty("device")
        String deviceId
) {
}
