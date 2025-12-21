package dev.idank.whisper.requests;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Set;

@JsonIgnoreProperties(ignoreUnknown = false)
public record MessageRequest(
        @JsonProperty("message")
        String message,
        @JsonProperty("targets")
        Set<String> targets
) {

}
