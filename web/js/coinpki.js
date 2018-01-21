$(function(){
  $('#prove').click(function(){
    var walletaddr = $('#walletaddress').val();
    var wallettype = $('#wallettype').val();
    var prooftext = $('#prooftext').val();
    var sign = $('#sign').val();
    var backend = $('#backend').val();
    $('#prove').prop('disabled', true);
    var ws = new WebSocket(backend);
    ws.onopen = function(ev) {
      console.log('open');
      $('.status').hide();
      $('#connected').show();
      var msg = {"id": 0, "method": "prove", "params": []};
      msg.params[0] = walletaddr;
      msg.params[1] = prooftext;
      msg.params[2] = sign;
      ws.send(JSON.stringify(msg) + "\n");
    };
    ws.onclose = function(ev) {
      console.log('close');
      $('.status').hide();
      $('#disconnected').show();
    };
    ws.onmessage = function(ev) {
      console.log('message: ' + ev.data);
      $('.status').hide();
      $('#message').show();
      var json = JSON.parse(ev.data);
      var result = json.result;
      if (result) {
        console.log('result: ' + result);
      }
    };
    ws.onerror = function(ev) {
      console.log('error');
      $('.status').hide();
      $('#error').show();
    };
    return false;
  });
});
