function loginModal() {
   var element = document.getElementById("modal");
   element.classList.toggle("modal");
}

function loginSucces() {
   var element = document.getElementById("loginSucces");
   element.classList.remove("modal");
}
