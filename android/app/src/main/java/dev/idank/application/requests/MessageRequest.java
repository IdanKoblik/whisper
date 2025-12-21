package dev.idank.application.requests;

import com.fasterxml.jackson.annotation.JsonCreator;
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
