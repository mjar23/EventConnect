 // Get the user ID from the URL parameter
 var urlParams = new URLSearchParams(window.location.search);
 var userId = urlParams.get('userId');

 // Fetch the user's profile from the server
 var xhr = new XMLHttpRequest();
 xhr.open('GET', '/other-user-profile?userId=' + userId);

 xhr.onload = function() {
     if (xhr.status === 200) {
         var user = JSON.parse(xhr.responseText);
         // Display the user's profile details
         document.getElementById('username').textContent = user.username;
         document.getElementById('email').textContent = user.email;
         document.getElementById('firstName').textContent = user.firstName;
         document.getElementById('lastName').textContent = user.lastName;
         document.getElementById('bio').textContent = user.bio;
         document.getElementById('interests').textContent = user.interests;
         document.getElementById('location').textContent = user.location;
         document.getElementById('age').textContent = user.age;
         document.getElementById('instagramUsername').textContent = user.instagramUsername;
         document.getElementById('facebookUsername').textContent = user.facebookUsername;
         document.getElementById('snapchatUsername').textContent = user.snapchatUsername;

         // Get the "Add Friend" button element
         var addFriendBtn = document.getElementById('add-friend-btn');

         // Add click event listener
         addFriendBtn.addEventListener('click', function() {
             // Get the user ID from the URL parameter
             var friendId = urlParams.get('userId');

             // Get the current user's ID and token from localStorage or another source
             var userId = localStorage.getItem('userId');
           
             console.log('User ID:', userId);
             console.log('Token:', token);

             // Send the friend request to the server
             var xhr = new XMLHttpRequest();
             xhr.open('POST', '/send-friend-request');
             xhr.setRequestHeader('Content-Type', 'application/json');
             xhr.setRequestHeader('Authorization', 'Bearer ' + localStorage.getItem('token')); // Add this line

             xhr.onload = function() {
                 if (xhr.status === 200) {
                     console.log('Friend request sent successfully');
                     // Optionally, you can display a success message to the user
                 } else {
                     console.error('Error sending friend request:', xhr.status);
                     // Handle the error, such as displaying an error message to the user
                 }
             };

             xhr.onerror = function() {
                 console.error('Error sending friend request:', xhr.status);
                 // Handle the error, such as displaying an error message to the user
             };

             var requestData = {
                 UserID: userId,
                 FriendID: friendId
             };
             xhr.send(JSON.stringify(requestData));
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