var conn, reloadTimer, connectTimer;
    var connect = function() {
        conn = new WebSocket('ws://' + location.host + '/websocket');
        conn.onmessage = function(event) {
            if (event.data === 'change') {
                if (reloadTimer) clearTimeout(reloadTimer);
                reloadTimer = setTimeout(function() {
                    window.location.reload();
                }, 200);
            }
        };
        conn.onclose = function() {
            if (connectTimer) clearTimeout(connectTimer);
            connectTimer = setTimeout(function() {
                connect();
            }, 1000);
        };
    };
connect();