package com.example.mockusers.invoice;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.Instant;
import java.util.Map;

@RestController
@RequestMapping("/api")
public class InvoiceController {

    private final InvoiceStore store;

    public InvoiceController(InvoiceStore store) {
        this.store = store;
    }

    @PostMapping("/invoice")
    public ResponseEntity<Map<String, Object>> receiveInvoice(@RequestBody Map<String, Object> body) {
        store.add(new InvoiceEntry("POST", "/api/invoice", body, Instant.now()));
        return ResponseEntity.ok(Map.of("status", "OK", "body", body));
    }

    @PostMapping("/invoice-inspector/reset")
    public ResponseEntity<Map<String, String>> reset() {
        store.reset();
        return ResponseEntity.ok(Map.of("status", "reset"));
    }

    @GetMapping("/invoice-inspector/{email}")
    public ResponseEntity<InvoiceInspectorResponse> inspect(
            @PathVariable String email,
            @RequestParam(defaultValue = "0") long timeout) throws InterruptedException {

        long deadline = System.currentTimeMillis() + timeout;

        while (true) {
            var match = store.findLatestByEmail(email);
            if (match.isPresent()) {
                var e = match.get();
                return ResponseEntity.ok(new InvoiceInspectorResponse(
                        email, "fulfilled", e.receivedAt(), e.body()));
            }
            if (System.currentTimeMillis() >= deadline) {
                return ResponseEntity.ok(new InvoiceInspectorResponse(
                        email, "waiting", null, null));
            }
            Thread.sleep(200);
        }
    }
}
