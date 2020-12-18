var global = {
  x: undefined,
  y: undefined,
  clickAttr: {
    click: false,
    side: "left",
  },
};
var word = "";
function writeKey() {
  word += event.key;
}
function clicki(e) {
  if (event.button == 2) {
    global.clickAttr = {
      click: true,
      side: "rigth",
    };
  } else {
    global.clickAttr = {
      click: true,
      side: "left",
    };
  }
}
function SendCommand() {
  fetch(window.location.href + "command", {
    method: "POST",
    body: JSON.stringify({
      command: document.getElementById("command").value,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  document.getElementById("command").value = "";
}
function mouse_position() {
  let e = window.event;
  let img = document.getElementById("sc");
  global.x = Math.floor((e.clientX / img.width) * screen.width);
  global.y = Math.floor((e.clientY / img.height) * screen.height);
}
let time = setInterval(function () {
  let url = window.location.href + "image/{",
    random = Math.random();
  document.querySelector("#sc").src = `${url + random}}`;
  //--------------------->  this send the mouse position
  fetch(window.location.href + "mouse", {
    method: "POST",
    body: JSON.stringify(global),

    headers: {
      "Content-Type": "application/json",
    },
  });
  //---------------------> this send the text
  fetch(window.location.href + "typetext", {
    method: "POST",
    body: JSON.stringify({
      word: word,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  global.clickAttr = {
    click: false,
    side: "right",
  };
}, 500);
