const jsonServer = require("json-server");
const server = jsonServer.create();
const router = jsonServer.router("db.json");
const middlewares = jsonServer.defaults();

// In-memory stores
const receivedRequests = [];
const expectations = new Map();
let expectationCounter = 0;

function deepEqual(a, b) {
  return JSON.stringify(sortedKeys(a)) === JSON.stringify(sortedKeys(b));
}

function sortedKeys(obj) {
  if (typeof obj !== "object" || obj === null) return obj;
  if (Array.isArray(obj)) return obj.map(sortedKeys);
  return Object.keys(obj)
    .sort()
    .reduce((acc, k) => {
      acc[k] = sortedKeys(obj[k]);
      return acc;
    }, {});
}

server.use(middlewares);
server.use(jsonServer.bodyParser);

// POST /api/invoice
server.post("/api/invoice", (req, res) => {
  const body = req.body;
  const entry = {
    method: "POST",
    path: "/api/invoice",
    body,
    receivedAt: new Date(),
  };
  receivedRequests.push(entry);

  for (const [, expectation] of expectations) {
    if (
      expectation.method === "POST" &&
      expectation.path === "/api/invoice" &&
      deepEqual(expectation.body, body)
    ) {
      expectation.fulfilledAt = new Date();
    }
  }

  res.status(200).json({ status: "OK", body });
});

// POST /api/expect
server.post("/api/expect", (req, res) => {
  const { method, path, body, id } = req.body;
  const uniqueId = String(id || ++expectationCounter);
  expectations.set(uniqueId, { method, path, body, fulfilledAt: null });
  res
    .status(201)
    .json({
      id: uniqueId,
      status: "waiting",
      expectation: { method, path, body },
    });
});

// GET /api/expect/:id
server.get("/api/expect/:id", async (req, res) => {
  const { id } = req.params;
  const timeout = parseInt(req.query.timeout) || 0;

  if (!expectations.has(id)) {
    return res.status(404).json({ error: "Expectation not found" });
  }

  const deadline = Date.now() + timeout;

  const check = () => {
    const expectation = expectations.get(id);
    if (!expectation.fulfilledAt) {
      const match = receivedRequests.find(
        (r) =>
          r.method === expectation.method &&
          r.path === expectation.path &&
          deepEqual(r.body, expectation.body),
      );
      if (match) {
        expectation.fulfilledAt = match.receivedAt;
      }
    }
    if (expectation.fulfilledAt) {
      return res
        .status(200)
        .json({
          id,
          status: "fulfilled",
          fulfilledAt: expectation.fulfilledAt,
        });
    }
    if (Date.now() >= deadline) {
      return res.status(200).json({ id, status: "waiting" });
    }
    setTimeout(check, 200);
  };

  check();
});

// DELETE /api/expect/:id
server.delete("/api/expect/:id", (req, res) => {
  expectations.delete(req.params.id);
  res.status(204).send();
});

// Add your custom route rewrite here
server.use(
  jsonServer.rewriter({
    "/api/*": "/$1",
  }),
);

server.use(router);

const PORT = process.env.PORT || 9000;

server.listen(PORT, () => {
  console.log(`JSON Server is running on http://localhost:${PORT}`);
});
