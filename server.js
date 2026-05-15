const jsonServer = require("json-server");
const server = jsonServer.create();
const router = jsonServer.router("db.json");
const middlewares = jsonServer.defaults();

// In-memory stores
const receivedRequests = [];

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

  res.status(200).json({ status: "OK", body });
});

// POST /api/invoice-inspector/reset
server.post("/api/invoice-inspector/reset", (req, res) => {
  receivedRequests.length = 0;
  res.status(200).json({ status: "reset" });
});

// GET /api/invoice-inspector/:email
server.get("/api/invoice-inspector/:email", (req, res) => {
  const email = req.params.email;

  const check = () => {
    const match = [...receivedRequests]
      .reverse()
      .find((r) => r.body && r.body.email === email);

    if (match) {
      return res.status(200).json({
        email,
        status: "fulfilled",
        fulfilledAt: match.receivedAt,
        body: match.body,
      });
    }
    return res.status(200).json({ email, status: "waiting" });
  };

  check();
});

server.post("/api/invoice-inspector", (req, res) => {
  const email = req.body.email;

  const check = () => {
    const match = [...receivedRequests]
      .reverse()
      .find((r) => r.body && r.body.email === email);

    if (match) {
      return res.status(200).json({
        email,
        status: "fulfilled",
        fulfilledAt: match.receivedAt,
        body: match.body,
      });
    }
    return res.status(200).json({ email, status: "waiting" });
  };

  check();
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
