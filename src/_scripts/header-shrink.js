document.addEventListener('DOMContentLoaded', function() {
    const header = document.querySelector('.page-header');
    if (!header) return;
    
    let isScrolled = false;
    const headerContainer = document.querySelector('.header-container');
    
    function handleScroll() {
        // Use container height for consistent calculation regardless of header state
        const fullHeight = headerContainer.offsetHeight;
        const triggerPoint = fullHeight - 60;
        const scrolled = window.scrollY > triggerPoint;
        
        if (scrolled !== isScrolled) {
            isScrolled = scrolled;
            header.classList.toggle('shrunk', scrolled);
        }
    }
    
    window.addEventListener('scroll', handleScroll);
    handleScroll(); // Check initial state
});