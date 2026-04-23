package com.example.mockusers.invoice;

import java.time.Instant;

public record InvoiceEntry(
    String method,
    String path,
    Object body,
    Instant receivedAt
) {}
