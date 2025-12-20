package dev.idank.application.clients;

import android.content.Context;
import android.util.Log;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import dev.idank.application.listeners.WebsocketListener;
import lombok.Getter;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.WebSocket;

public class WebsocketClient {

    private static WebsocketClient instance;

    private final OkHttpClient client;
    private final String apiToken;
    private final String websocketURL;
    private final WebsocketListener wsListener;

    @Getter
    private final String deviceID;

    private WebSocket websocket;

    private volatile boolean connecting = false;
    private volatile boolean connected = false;

    public static synchronized void setInstance(String apiToken, String websocketURL, String deviceID) {
        if (instance == null) {
            instance = new WebsocketClient(apiToken, websocketURL, deviceID);
        }
    }

    public static WebsocketClient getInstance() {
        if (instance == null) {
            throw new IllegalStateException("WebsocketClient not initialized");
        }
        return instance;
    }

    private WebsocketClient(String apiToken, String websocketURL, String deviceID) {
        this.apiToken = apiToken;
        this.websocketURL = websocketURL;
        this.deviceID = deviceID;

        this.client = new OkHttpClient.Builder()
                .readTimeout(0, TimeUnit.MILLISECONDS)
                .build();

        this.wsListener = new WebsocketListener(this);
    }

    public synchronized void connect() {
        if (connecting || connected) return;

        connecting = true;

        Request request = new Request.Builder()
                .url(websocketURL)
                .addHeader("X-Api-Token", apiToken)
                .build();

        websocket = client.newWebSocket(request, wsListener);

        Log.d(getClass().getName(), "Connecting to " + websocketURL);
    }

    public CompletableFuture<Void> disconnect() {
        CompletableFuture<Void> future = new CompletableFuture<>();

        if (websocket == null) {
            future.complete(null);
            return future;
        }

        boolean closed = websocket.close(1000, "Client closed");

        if (closed) {
            future.complete(null);
        } else {
            future.completeExceptionally(new RuntimeException("Failed to close websocket"));
        }

        return future;
    }

    public synchronized void reconnect() {
        if (connected || connecting) return;

        disconnect().whenComplete((v, e) -> connect());
    }

    public synchronized void shutdown() {
        disconnect().whenComplete((v, e) -> {
            connected = false;
            connecting = false;
            instance = null;
            Log.d(getClass().getName(), "WebSocket shutdown");
        });
    }

    public synchronized void onConnected() {
        connected = true;
        connecting = false;
    }

    public synchronized void onDisconnected() {
        connected = false;
        connecting = false;
    }
}
