document.addEventListener('DOMContentLoaded', function() {
    const navButtons = document.querySelectorAll('.corner-nav');
    const bizCard = document.querySelector('.bizCard');
    
    navButtons.forEach(button => {
        button.addEventListener('mouseenter', function() {
            bizCard.classList.add('nav-hovered');
        });
        
        button.addEventListener('mouseleave', function() {
            bizCard.classList.remove('nav-hovered');
        });
    });
});