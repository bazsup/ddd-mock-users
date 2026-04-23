package com.example.mockusers.invoice;

import org.springframework.stereotype.Component;

import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CopyOnWriteArrayList;

@Component
public class InvoiceStore {

    private final CopyOnWriteArrayList<InvoiceEntry> entries = new CopyOnWriteArrayList<>();

    public void add(InvoiceEntry entry) {
        entries.add(entry);
    }

    public void reset() {
        entries.clear();
    }

    public Optional<InvoiceEntry> findLatestByEmail(String email) {
        var snapshot = entries;
        for (int i = snapshot.size() - 1; i >= 0; i--) {
            var e = snapshot.get(i);
            if (e.body() instanceof Map<?, ?> map && email.equals(map.get("email"))) {
                return Optional.of(e);
            }
        }
        return Optional.empty();
    }
}
