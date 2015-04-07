var out = document.getElementById('console');

function log(format, var_args) {
  for (var i = 1; i < arguments.length; i++) {
    format = format.replace(/%s/, JSON.stringify(arguments[i]));
  }
  var span = document.createElement('span');
  span.textContent = format + '\n';
  out.appendChild(span);
}

var ws = new WebSocket('ws://localhost:1987/socket');

function send(type, value) {
  var msg = {Type: type, Value: value};
  log('> %s', msg);
  ws.send(JSON.stringify(msg));
}

ws.onclose = function () {
  log('Connection closed.');
};

ws.onopen = function () {
  log('Connected.');
  send('Hello', {SessionId: '', ClientVersion: '0.1.0.001'});
};

ws.onmessage = function (event) {
  var msg = JSON.parse(event.data);
  log('< %s', msg);
  switch (msg.Type) {
    case 'Ping':
      send('Ping', {Id: msg.Value.Id, Time: +new Date()});
      break;
  }
};
