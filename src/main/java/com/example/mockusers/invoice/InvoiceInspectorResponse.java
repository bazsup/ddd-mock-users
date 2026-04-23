package com.example.mockusers.invoice;

import com.fasterxml.jackson.annotation.JsonInclude;
import java.time.Instant;

@JsonInclude(JsonInclude.Include.NON_NULL)
public record InvoiceInspectorResponse(
    String email,
    String status,
    Instant fulfilledAt,
    Object body
) {}
