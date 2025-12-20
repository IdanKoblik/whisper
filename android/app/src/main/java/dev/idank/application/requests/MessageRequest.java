package dev.idank.application.requests;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Set;

public record MessageRequest(
        String device,
        String message,
        Set<String> targets
) {

}
