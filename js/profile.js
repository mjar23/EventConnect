// Check if the user is logged in
var token = localStorage.getItem('token');
if (!token) {
    // Redirect to the login page if the user is not logged in
    window.location.href = '/login.html';
}

// Fetch user profile and activities from the server
var xhr = new XMLHttpRequest();
xhr.open('GET', '/profile');
xhr.setRequestHeader('Authorization', 'Bearer ' + token);

xhr.onload = function() {
    if (xhr.status === 200) {
        var data = JSON.parse(xhr.responseText);
        console.log('Response data:', data);
        console.log('User profile:', data.user);
        console.log('User activities:', data.activities);
        console.log('User teams:', data.teams);

        // Display user profile details
        document.getElementById('username').textContent = data.user.username || '';
        document.getElementById('email').textContent = data.user.email || '';
        document.getElementById('firstName').textContent = data.user.firstName || '';
        document.getElementById('lastName').textContent = data.user.lastName || '';
        document.getElementById('bio').textContent = data.user.bio || '';
        document.getElementById('interests').textContent = data.user.interests || '';
        document.getElementById('location').textContent = data.user.location || '';
        document.getElementById('age').textContent = data.user.age || '';
        // Update to display social media usernames
        document.getElementById('instagramUsername').textContent = data.user.instagramUsername || '';
        document.getElementById('facebookUsername').textContent = data.user.facebookUsername || '';
        document.getElementById('snapchatUsername').textContent = data.user.snapchatUsername || '';

        // Display user events
        const eventRegistrationList = document.getElementById('event-registration-list');
        const loadMoreBtn = document.getElementById('load-more-btn');

        function fetchEventDetails(eventId) {
            const apiKey = "3da8b2c4d34b51e8a45e86cc08280ad9";
            const url = `https://www.skiddle.com/api/v1/events/${eventId}/?api_key=${apiKey}`;

            return fetch(url)
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error! status: ${response.status}`);
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.results) {
                        return data.results;
                    } else {
                        throw new Error('No results in the response');
                    }
                })
                .catch(error => {
                    console.error('Error fetching event details:', error);
                    return null;
                });
        }

        function createEventCard(eventDetails) {
            const eventCard = document.createElement('div');
            eventCard.classList.add('event-card');

            const eventCardHeader = document.createElement('div');
            eventCardHeader.classList.add('event-card-header');
            const eventImageElement = document.createElement('img');
            eventImageElement.src = eventDetails.imageurl || '/path/to/default/image.jpg';
            eventImageElement.alt = eventDetails.eventname || 'Event';
            const eventTitleElement = document.createElement('div');
            eventTitleElement.classList.add('event-title');
            eventTitleElement.textContent = eventDetails.eventname || 'Unnamed Event';
            eventCardHeader.appendChild(eventImageElement);
            eventCardHeader.appendChild(eventTitleElement);

            const eventCardContent = document.createElement('div');
            eventCardContent.classList.add('event-card-content');
            const eventDateElement = document.createElement('p');
            eventDateElement.textContent = eventDetails.date ? `Date: ${eventDetails.date}` : 'Date: Not specified';
            const eventLinkElement = document.createElement('a');
            eventLinkElement.href = `/event-details.html?eventId=${eventDetails.id}`;
            eventLinkElement.textContent = 'View Details';
            eventCardContent.appendChild(eventDateElement);
            eventCardContent.appendChild(eventLinkElement);

            eventCard.appendChild(eventCardHeader);
            eventCard.appendChild(eventCardContent);

            return eventCard;
        }

        if (data.activities && data.activities.length > 0) {
            console.log('Activities:', data.activities);
            let shownEvents = 0;
            data.activities.forEach(function(activity) {
                if (activity.activity_type === 'event_registered') {
                    console.log('Registered event:', activity);
                    if (shownEvents < 5) { // Display only the first 5 events
                        fetchEventDetails(activity.event_id)
                            .then(eventDetails => {
                                console.log('Fetched event details:', eventDetails);
                                if (eventDetails) {
                                    const eventCard = createEventCard(eventDetails);
                                    eventRegistrationList.appendChild(eventCard);
                                    shownEvents++;
                                }
                            })
                            .catch(error => console.error('Error fetching event details:', error));
                    }
                }
            });

            if (data.activities.length > 5) {
                loadMoreBtn.style.display = 'block';
                loadMoreBtn.addEventListener('click', function() {
                    // Display the remaining events
                    for (let i = 5; i < data.activities.length; i++) {
                        if (data.activities[i].activity_type === 'event_registered') {
                            fetchEventDetails(data.activities[i].event_id)
                                .then(eventDetails => {
                                    if (eventDetails) {
                                        const eventCard = createEventCard(eventDetails);
                                        eventRegistrationList.appendChild(eventCard);
                                    }
                                })
                                .catch(error => console.error('Error fetching event details:', error));
                        }
                    }
                    loadMoreBtn.style.display = 'none';
                });
            }
        } else {
            const noEventsMessage = document.createElement('div');
            noEventsMessage.classList.add('event-card');
            noEventsMessage.innerHTML = '<div class="event-card-header"><div class="event-title">No Events</div></div><div class="event-card-content"><p>You have not registered for any events yet.</p></div>';
            eventRegistrationList.appendChild(noEventsMessage);
        }

        // Display user teams
        console.log('Data received from server:', data);

        // Add event listeners for edit profile functionality
        var editProfileBtn = document.getElementById('edit-profile-btn');
        var editProfileForm = document.getElementById('edit-profile-form');
        var profileDetails = document.getElementById('profile-details');
        var updateProfileForm = document.getElementById('update-profile-form');
        var cancelEditBtn = document.getElementById('cancel-edit-btn');

        console.log('Edit Profile button:', editProfileBtn); // Debug output

        editProfileBtn.addEventListener('click', function() {
            // Show the edit profile form and hide the profile details
            editProfileForm.style.display = 'block';
            profileDetails.style.display = 'none';

            // Populate the form fields with the current user data
            document.getElementById('edit-username').value = data.user.username;
            document.getElementById('edit-email').value = data.user.email;
            document.getElementById('edit-firstName').value = data.user.firstName;
            document.getElementById('edit-lastName').value = data.user.lastName;
            document.getElementById('edit-bio').value = data.user.bio;
            document.getElementById('edit-interests').value = data.user.interests;
            document.getElementById('edit-location').value = data.user.location;
            document.getElementById('edit-instagramUsername').value = data.user.instagramUsername;
            document.getElementById('edit-facebookUsername').value = data.user.facebookUsername;
            document.getElementById('edit-snapchatUsername').value = data.user.snapchatUsername;
        });

        cancelEditBtn.addEventListener('click', function() {
            // Hide the edit profile form and show the profile details
            editProfileForm.style.display = 'none';
            profileDetails.style.display = 'block';
        });

        updateProfileForm.addEventListener('submit', function(event) {
            event.preventDefault();

            // Get the updated profile data from the form
            var updatedProfile = {
                username: document.getElementById('edit-username').value,
                email: document.getElementById('edit-email').value,
                firstName: document.getElementById('edit-firstName').value,
                lastName: document.getElementById('edit-lastName').value,
                bio: document.getElementById('edit-bio').value,
                interests: document.getElementById('edit-interests').value,
                location: document.getElementById('edit-location').value,
                age: data.user.age, // Include the current age value
                instagramUsername: document.getElementById('edit-instagramUsername').value,
                facebookUsername: document.getElementById('edit-facebookUsername').value,
                snapchatUsername: document.getElementById('edit-snapchatUsername').value
            };

            // Make an HTTP request to update the user profile
            var xhr = new XMLHttpRequest();
            xhr.open('PUT', '/profile');
            xhr.setRequestHeader('Content-Type', 'application/json');
            xhr.setRequestHeader('Authorization', 'Bearer ' + token);

            xhr.onload = function() {
                if (xhr.status === 200) {
                    // Profile updated successfully
                    alert('Profile updated successfully');
                    // Refresh the profile details on the page
                    location.reload();
                } else {
                    // Handle error
                    alert('Error updating profile');
                }
            };

            xhr.send(JSON.stringify(updatedProfile));
        });
    } else {
        console.error('Error fetching user profile:', xhr.status);
        // Handle the error, such as displaying an error message to the user
    }
};

xhr.onerror = function() {
    console.error('Error fetching user profile:', xhr.status);
    // Handle the error, such as displaying an error message to the user
};

xhr.send();