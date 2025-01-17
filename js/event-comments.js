// Function to load event details
async function loadEventDetails() {
    const urlParams = new URLSearchParams(window.location.search);
    const eventId = urlParams.get('eventId');

    try {
        const response = await fetch(`http://localhost:8000/events/${eventId}`);
        if (!response.ok) {
            throw new Error('Error loading event details');
        }
        const event = await response.json();
        displayEventDetails(event);
    } catch (error) {
        console.error('Error loading event details:', error);
    }
}

// Function to display event details
function displayEventDetails(event) {
    const eventImage = document.getElementById('event-image');
    const eventName = document.getElementById('event-name');
    const eventDate = document.getElementById('event-date');
    const eventDescription = document.getElementById('event-description');

    eventImage.src = event.imageurl;
    eventName.textContent = event.eventname;
    eventDate.textContent = `Date: ${event.date}`;
    eventDescription.textContent = event.description;
}

// Function to load comments for the event
async function loadComments() {
    const urlParams = new URLSearchParams(window.location.search);
    const eventId = urlParams.get('eventId');

    try {
        const response = await fetch(`http://localhost:8000/events/${eventId}/comments`);
        const comments = await response.json();
        displayComments(comments);
    } catch (error) {
        console.error('Error loading comments:', error);
    }
}

// Function to display comments
function displayComments(comments) {
    const commentsContainer = document.getElementById('comments-container');
    commentsContainer.innerHTML = '';

    comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.classList.add('comment');
        commentElement.innerHTML = `
            <p><a href="#" onclick="viewUserProfile(${comment.userId})">${comment.username}</a> - ${comment.createdAt}</p>
            <p>${comment.text}</p>
        `;
        commentsContainer.appendChild(commentElement);
    });
}

// Function to view the other user's profile
async function viewUserProfile(userId) {
    try {
        const response = await fetch(`http://localhost:8000/other-user-profile?userId=${userId}`);
        if (!response.ok) {
            throw new Error('Error loading user profile');
        }
        const userProfile = await response.json();
        displayOtherUserProfile(userProfile);
    } catch (error) {
        console.error('Error loading user profile:', error);
    }
}

// Function to display the other user's profile
function displayOtherUserProfile(userProfile) {
    // Redirect to the other-user-profile.html page with the user profile data
    window.location.href = `/other-user-profile.html?userId=${userProfile.id}`;
}

// Function to submit a new comment
async function submitComment() {
    const urlParams = new URLSearchParams(window.location.search);
    const eventId = urlParams.get('eventId');
    const commentInput = document.getElementById('comment-input');
    const commentText = commentInput.value.trim();

    if (commentText === '') {
        return;
    }

    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`http://localhost:8000/events/${eventId}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ text: commentText })
        });

        if (response.ok) {
            commentInput.value = '';
            loadComments();
        } else {
            console.error('Error submitting comment:', response.status);
        }
    } catch (error) {
        console.error('Error submitting comment:', error);
    }
}

// Load event details and comments on page load
loadEventDetails();
loadComments();
