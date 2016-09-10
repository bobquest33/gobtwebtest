function validateTerms(){
  var c=document.getElementById('termsCheckbox');
  var d=document.getElementById('terms_div');
  if (c.checked) {
    d.innerHTML ="On";
	sendRequest();
    return true;
  } else { 
    d.innerHTML = "Off";
	sendRequest();


    return false;
  }
}


this.sendRequest = function(){
   var url="/api";
   var d=document.getElementById('terms_div');
   var data= d.innerHTML;
   var client = new XMLHttpRequest();
   client.open("POST", url, true);
   client.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
   client.setRequestHeader("Connection", "close");
   client.send("switch=" + encodeURIComponent(data));   


   if (client.status == 200){
      d.innerHTML = client.responseText;
   }
   else{
       d.innerHTML = client.statusText;
   }

  }

