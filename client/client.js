var jayson = require('jayson');

// create a client
var client = jayson.client.http({
  port: 9090
});

client.request('initialize', [{"processId": 1, "rootPath": "/", "capabilities": {}}], function(err, response) {
  if (err) throw err;
  console.log(response.result);
});
