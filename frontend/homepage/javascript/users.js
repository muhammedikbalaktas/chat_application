
const profileInfoElements = document.getElementsByClassName('users_div');


for (let i = 0; i < profileInfoElements.length; i++) {
    profileInfoElements[i].addEventListener('click', function() {
        
        alert('Clicked on profile info');
    });
}
