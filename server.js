const jsonServer = require('json-server');
const server = jsonServer.create();
const router = jsonServer.router('db.json');
const middlewares = jsonServer.defaults();

server.use(middlewares);

// Add your custom route rewrite here
server.use(jsonServer.rewriter({
  '/api/*': '/$1',
}));

server.use(router);

const PORT = process.env.PORT || 8080

server.listen(PORT, () => {
  console.log(`JSON Server is running on http://localhost:${PORT}`);
});