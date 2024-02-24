function validatePassword(){
  const password = document.getElementById('password');
  const confirm_password = document.getElementById('repeat-psw');
  if(password.value != confirm_password.value) {
    confirm_password.setCustomValidity("Passwords Don't Match");
  } else {
    confirm_password.setCustomValidity("");
  }

  confirm_password.reportValidity();
}

window.addEventListener('DOMContentLoaded', () => {
  const password = document.getElementById('password');
  const confirm_password = document.getElementById('repeat-psw');

  password.onchange = validatePassword;
  confirm_password.onkeyup = validatePassword;
})