// Get the event ID from the URL query parameter
const urlParams = new URLSearchParams(window.location.search);
const eventId = urlParams.get('eventId');

let userLocations = [];
let map;

async function getEventDetails(eventId) {
    try {
        const [eventResponse, userLocationsResponse] = await Promise.all([
            fetch(`http://localhost:8000/events/${eventId}`),
            fetch(`http://localhost:8000/events/${eventId}/user-locations`)
        ]);

        const eventData = await eventResponse.json();
        userLocations = await userLocationsResponse.json();

        console.log('User Locations:', userLocations);

        renderEventDetails(eventData);
        loadGoogleMapsAPI();

        // Pass the eventData object to getTwitterData
        getTwitterData(eventData);
    } catch (error) {
        console.error('Error fetching event details:', error);
    }
}

function renderEventDetails(event) {
    const eventContainer = document.getElementById('event-container');

    eventContainer.innerHTML = `
        <div class="event-header">
            <div class="event-image">
                <img src="${event.imageurl}" alt="${event.eventname}">
            </div>
            <h2 class="event-title">${event.eventname}</h2>
        </div>
        <div class="event-content">
            <div class="event-details">
                <div class="detail">
                    <h3><i class="fas fa-info-circle"></i> Description</h3>
                    <p>${event.description}</p>
                </div>
                <div class="detail">
                    <h3><i class="fas fa-calendar-alt"></i> Date</h3>
                    <p>${event.date}</p>
                </div>
                <div class="detail">
                    <h3><i class="fas fa-map-marker-alt"></i> Venue</h3>
                    <p>${event.venue.name}</p>
                </div>
                <div class="detail">
                    <h3><i class="fas fa-dollar-sign"></i> Entry Price</h3>
                    <p>${event.entryprice}</p>
                </div>
                <div class="detail">
                    <h3><i class="fas fa-user"></i> Minimum Age</h3>
                    <p>${event.minage}</p>
                </div>
                <div class="detail">
                    <h3><i class="fas fa-link"></i> TicketLink</h3>
                    <p><a href="${event.link}" target="_blank">${event.link}</a></p>
                </div>
            </div>
        </div>
        <div class="event-actions">
            <button class="action-button" onclick="registerEvent(${event.id})"><i class="fas fa-heart"></i> I'm Interested</button>
            <button class="action-button" onclick="openCommentForum(${event.id})"><i class="fas fa-comments"></i> Comment Forum</button>
            <button class="action-button" onclick="enterRaffle(${event.id})"><i class="fas fa-ticket-alt"></i> Enter Raffle</button>
        </div>
    `;
}

function initMap() {
    if (userLocations.length === 0) {
        // Set default center and zoom level if no user locations available
        map = new google.maps.Map(document.getElementById('map'), {
            center: { lat: 0, lng: 0 },
            zoom: 2,
        });
        return;
    }

    const bounds = new google.maps.LatLngBounds();
    userLocations.forEach(location => {
        const latLng = new google.maps.LatLng(location.latitude, location.longitude);
        bounds.extend(latLng);
    });

    map = new google.maps.Map(document.getElementById('map'), {
        center: bounds.getCenter(),
        zoom: 8,
    });

    map.fitBounds(bounds);

    plotUserLocations(userLocations);
}

function plotUserLocations(userLocations) {
    userLocations.forEach(location => {
        const latLng = new google.maps.LatLng(location.latitude, location.longitude);
        const marker = new google.maps.Marker({
            position: latLng,
            map: map,
            userId: location.id // Store the user ID as a data attribute
        });

        console.log('User ID:', location.id);

        // Add click event listener to the marker
        marker.addListener('click', () => {
            // Redirect to other-user-profile.html with the user ID as a query parameter
            const userId = marker.get('userId');
            window.location.href = `/other-user-profile.html?userId=${userId}`;
        });
    });
}

function loadGoogleMapsAPI() {
    const script = document.createElement('script');
    script.src = `https://maps.googleapis.com/maps/api/js?key=AIzaSyBDR_N5soK76D55aICV-bW2jDEOxE5WwVQ&callback=initMap`;
    script.async = true;
    script.defer = true;
    document.head.appendChild(script);
}

async function registerEvent(eventId) {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`http://localhost:8000/events/${eventId}/register`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (response.ok) {
            console.log('Event registered successfully');
            // Update the UI or show a success message
        } else {
            console.error('Error registering event:', response.status);
            // Show an error message to the user
        }
    } catch (error) {
        console.error('Error registering event:', error);
        // Show an error message to the user
    }
}

async function enterRaffle(eventId) {
    try {
        const token = localStorage.getItem('token');
        const userInfo = await getUserInfo();

        const response = await fetch(`http://localhost:8000/events/${eventId}/raffle`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                eventId: eventId.toString(),
                age: userInfo.age,
                gender: userInfo.gender,
                latitude: userInfo.latitude,
                longitude: userInfo.longitude
            })
        });

        if (response.ok) {
            console.log('Entered raffle successfully');
            // Update the UI or show a success message
        } else {
            console.error('Error entering raffle:', response.status);
            // Show an error message to the user
        }
    } catch (error) {
        console.error('Error entering raffle:', error);
        // Show an error message to the user
    }
}

async function getUserInfo() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('http://localhost:8000/user', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (response.ok) {
            const userInfo = await response.json();
            return userInfo;
        } else {
            console.error('Error fetching user info:', response.status);
            throw new Error('Failed to fetch user info');
        }
    } catch (error) {
        console.error('Error fetching user info:', error);
        throw error;
    }
}

function openCommentForum(eventId) {
    window.location.href = `/event-comments.html?eventId=${eventId}`;
}

async function getTwitterData(event) {
    try {
        const eventName = event.eventname; // Extract the event name from the event data
        const response = await fetch(`/events/${encodeURIComponent(eventName)}/twitter-scraper`);
        const twitterData = await response.json();
        renderTwitterData(twitterData);
    } catch (error) {
        console.error('Error fetching Twitter data:', error);
        // Handle the error, such as displaying an error message to the user
    }
}

function renderTwitterData(twitterData) {
    const tweetContainer = document.getElementById('tweet-container');

    console.log('Twitter data:', twitterData);

    if (!tweetContainer) {
        console.error('Tweet container not found');
        return;
    }

    tweetContainer.innerHTML = '';

    if (twitterData.length === 0) {
        tweetContainer.innerHTML = '<p>No tweets found.</p>';
        return;
    }

    const tweetBoxTitle = document.createElement('h3');
    tweetBoxTitle.innerHTML = '<i class="fab fa-twitter"></i> Related Tweets';
    tweetContainer.appendChild(tweetBoxTitle);

    twitterData.forEach(tweet => {
        const tweetElement = document.createElement('div');
        tweetElement.classList.add('tweet-card');
        tweetElement.innerHTML = `
            <p class="tweet-text">${tweet.text}</p>
            <p class="tweet-author">By: ${tweet.user.name} (@${tweet.user.screen_name})</p>
        `;
        tweetContainer.appendChild(tweetElement);
    });
}
function goBack() {
    window.location.href = '/events.html';
}

// Fetch and render the event details
getEventDetails(eventId);