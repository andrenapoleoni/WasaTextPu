<script>
export default {
    data: function(){
        return{
            userId: localStorage.userId,
            username: localStorage.username,
            photo: localStorage.photo,
            newUsername: "",
            newPhoto: null,
            showBigPhoto: false,
        };

    },

    methods: {
        async updateUsername()
        {
            console.log(this.newUsername);
            try{
                let response=await this.$axios.put(
                     `/user/${sessionStorage.userID}/username`,
                     {username: this.newUsername},
                     {headers: {Authorization: sessionStorage.token}}

                );
                this.username = this.newUsername;
                sessionStorage.username = this.newUsername;
                alert(response.data);
                this.newUsername = "";
                
                
            }catch(error){
              console.log(error);

              // Controlla se l'errore è un Internal Server Error (500)
              if (error.response && error.response.status === 500) {
                // Mostra un messaggio specifico per l'errore
                alert("Username already taken");
              } else {
                // Gestione di altri errori
                alert("An unexpected error occurred. Please try again.");
              }
              this.newUsername = "";

                
                
                
            }
        },
        
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    changePhoto(event) {
      const file = event.target.files[0];
      if (file) {
        this.newPhoto = file;
        this.updatePhoto();
      }
    },
    async updatePhoto() {
      const formData = new FormData();
      formData.append('image', this.newPhoto);
      try {
        const response = await this.$axios.put(
          `/user/${sessionStorage.userID}/photo`,
          formData,
          {
            headers: {
              Authorization: sessionStorage.token,
            },
          }
        );
        this.photo = response.data.photo;
        sessionStorage.photo = response.data.photo;
        this.newPhoto = null;
        // Ricarica la pagina dopo aver aggiornato la foto
       
      } catch (error) {
        console.log(error);
      }
    },
    showphoto(){
      this.showBigPhoto=!this.showBigPhoto;
    }

},
};




</script>
<template> 
<div class="user-info">
    <img :src="`data:image/jpg;base64,${photo}`" alt="Profile Picture" class="user-photo" @click="showphoto"/>
    <h2>{{ username }}</h2>
    <div class="change-username">
      <input
        v-model="newUsername"
        type="text"
        placeholder="Enter new username"
      />
      <button @click="updateUsername">Change Username</button>
    </div>
    <div class="change-photo">
      <input type="file" ref="fileInput" @change="changePhoto" style="display: none;" />
      <button @click="triggerFileInput">Change photo</button>
    </div>
  </div>
  <div class="bigphoto" v-if="showBigPhoto">
    <img :src="`data:image/jpg;base64,${photo}`" alt="Profile Picture" class="user-bigphoto" @click="showphoto"/>
  </div>

</template>

<style>
.bigphoto {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5); /* Semi-transparent background */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(10px); /* Apply the blur effect */
}

.user-bigphoto {
  width: 70vh; /* Use 70% of the viewport height */
  height: 70vh;
  max-width: 80%;
  max-height: 80%;
  border-radius: 50%;
  object-fit: cover; /* Ensure the image covers the container proportionally */
}
</style>


