
<script>

import Reaction from "../components/ReactionModal.vue";
import Forward from "../components/ForwardModal.vue";

export default{ 
  data() {
    return {
      owner: sessionStorage.username,
      ownerID: sessionStorage.userID,
      messages: [],
      newMessage: "",
      convID: "null" , // Handle new chat
      userToSend: localStorage.username,
      userIdToSend: localStorage.userID,
      avatar: localStorage.photo,
      groupname : localStorage.groupname,
      groupId : localStorage.groupId,
      groupMembers: JSON.parse(localStorage.groupMembers || "[]"),
      usertoadd: [],
      showModal: false,
      searchTerm: "",
      searchResults: [],
      dropdownIndex: null,
      modalReaction: false,
      selectedMessage: null,
      modalForward: false,
      convList: JSON.parse(localStorage.conversationList || "[]"),
      intervalId: null,
      showBigPhoto:false,
      photoToSend:null,
      errormsg: null,
      showphotomessage : false,
      selectedMessageph : null,
      replyId: null,
      comments: [],

    };
  },
  methods: {
    async sendMessage() {
        if (this.newMessage=="" && this.photoToSend== null) throw "impossible send an empty message"
      
        try {
          // Send the message to the server
          console.log("ConvID:", this.convID,"groupID:",this.groupId,"token:",sessionStorage.token,"userId:",sessionStorage.userID);
          const formData = new FormData();
          if (this.photoToSend) {
            formData.append("file", this.photoToSend);
          }else{
            formData.append("file", null);
          }
          
            formData.append("text", this.newMessage);
          if(this.replyId){
            formData.append("reply", this.replyId);
          }else{
            formData.append("reply", "");
          }  
            
            await this.$axios.post(
              `/user/${sessionStorage.userID}/conversation/${this.convID}/messages`,
              formData,
              { headers: { Authorization: sessionStorage.token } }
            );
          
          console.log("Message sent:", this.newMessage);
          this.photoToSend = null; // Reset the photo to send
          
          
          

          this.newMessage = ""; // Clear the input field
          this.replyId = null; // Reset the reply ID
          this.getConversation();
        } catch (e) {
          alert("Error sending message: " + e.toString());
        }
      
    },

    async createConversation() {
      try {
        const response = await this.$axios.put(
          `/user/${sessionStorage.userID}/conversation/${this.userIdToSend}`,
          {txt: this.newMessage},
          { headers: { Authorization: sessionStorage.token } }
        );
        this.convID = response.data.conversationID; // Update convID with the backend value
        
        

        this.newMessage = ""; // Clear the input field
        
        
        this.$router.push(`/conversation/${this.convID}`);
        this.getConversation();
      } catch (e) {
        alert("Error creating conversation: " + e.toString());
        throw e;
      }
    },
    showphoto(){
      this.showBigPhoto=!this.showBigPhoto;
    },
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    sendPhoto(event) {
      const file = event.target.files[0];
      if (file) {
        this.photoToSend = file;
        console.log("File to send:", this.photoToSend);
        this.sendMessage();
        
      }
    },

    


    async getConversation() {
      try {
        if (this.convID !== "null") {
          const response = await this.$axios.get(
            `/user/${sessionStorage.userID}/conversation/${this.convID}`,
            { headers: { Authorization: sessionStorage.token } }
          );
          this.messages = response.data.messages;
          this.groupMembers = response.data.memberlist;
          console.log("Messages:", this.messages);
          console.log("username:", this.userToSend);
          console.log("convid", this.convID);
        }
      } catch (e) {
        alert("Error fetching conversation: " + e.toString());
      }
    },

    async addToGroup(){
      try{
        let response = await this.$axios.put(
          `/user/${sessionStorage.userID}/groups/${this.groupId}`,
          {users: this.usertoadd},
          { headers: { Authorization: sessionStorage.token } }
        );
        this.showModal = false;
        this.groupMembers = response.data.users;
       
        //aggiorna la conversazione
        this.getConversation();
      } catch (e) {
        alert("Error adding user to group: " + e.toString());

      }

    },
    
    
    async check() {
      // Check if the conversation ID is valid
      if (!this.convID || this.convID === "null") {
        try {
          await this.createConversation(); // Create the conversation if it doesn’t exist
        } catch (e) {
          alert("Error during conversation creation: " + e.toString());
          return;
        }
      } else {
        this.sendMessage(); // Send the message after ensuring the conversation exists
      }
    },

    toggleDropdown(index) {
      console.log("Index:", index);
      if (this.dropdownIndex === index) {
    this.dropdownIndex = null; // Chiude il dropdown se è già aperto
  } else {
    this.dropdownIndex = index; // Apre il dropdown per il messaggio cliccato
  }
  console.log("Dropdown index:", this.dropdownIndex);
    },
 
  

    opendModal()
    {
      this.showModal = true;
    },

    toggleModal()
    {
      this.showModal=false;
    },

    async searchUsers() {
      if (!this.searchTerm) {
        this.searchResults = [];
        return;
      }

      try {
        const response = await this.$axios.get(`/user/${sessionStorage.userID}`, {
          params: { search: this.searchTerm },
          headers: { 'Authorization': sessionStorage.token }
        });

        const groupMemberUsernames = this.groupMembers.map(member => member.username);
        this.searchResults = response.data.filter(user => !groupMemberUsernames.includes(user.username));
      } catch (e) {
        this.errormsg = e.toString();
      }
    },
    toggleUserSelection(user){
      const index = this.usertoadd.findIndex(u => u.username === user.username);
      if (index === -1) {
        this.usertoadd.push(user);
      } else {
        this.usertoadd.splice(index, 1);
      }

    },
    gotoGroupProfile(){
      if (Number(this.groupId) !== 0) {
        localStorage.groupId = this.groupId;
        localStorage.groupName = this.groupname;
        localStorage.groupMembers = JSON.stringify(this.groupMembers);
        localStorage.groupPhoto = this.avatar;
        this.$router.push('/group');
      }
    },
    async deleteMessage(msg) {
      console.log("ID del messaggio da eliminare:", msg.messageID, "msg:", msg);
      try {
        await this.$axios.delete(`/user/${sessionStorage.userID}/conversation/${this.convID}/messages/${msg.message.messageID}`, {
          headers: { 'Authorization': sessionStorage.token }
        });
        this.getConversation();
      } catch (e) {
        alert("Error deleting message: " + e.toString());
      }
    },

    async replyMessage(msg) {
      console.log("Replying to message:", msg);
      this.replyId = msg.message.messageID;
      
      
    },
    showModalReaction(msg) {
      console.log("msg", msg.message)
      this.selectedMessage= msg.message;
      this.comments= msg.comment;
      console.log("Selected message:", this.selectedMessage);
  
      this.modalReaction=!this.modalReaction;
    },
    showModalForward(msg) {
      console.log("lsit", this.convList)
      this.selectedMessage=msg.message;
      this.modalForward = !this.modalForward;

    },
    showphotomessageB(msg) {

      console.log("masg", msg)
      
      this.selectedMessageph= msg.message.photo;
      this.boolshowph();
      console.log("Selected message:", this.selectedMessageph);
      
  
      
    },
    boolshowph(){
      this.showphotomessage=!this.showphotomessage;
    },

  
    

    
  },
  
  mounted() {
    // Se l'utente non è loggato, reindirizza alla pagina di login
    if (!sessionStorage.token) {
      this.$router.push("/");
      return;
    }
    // Altrimenti ottente le conversazioni dell'utente
    this.convID = this.$route.params.conversation;
    console.log("Conversation ID:", this.convID);
    this.getConversation();
    this.intervalId = setInterval(async () => {
        clearInterval(this.intervalId);
        await this.getConversation();
        this.intervalId = setInterval(this.getConversation, 7000);
      }, 7000);
    },
    
  
  beforeunmount() {
    // ripulisci tutti i dati
    this.messages = [];
    this.newMessage = "";

    this.convID = "null";
  },

  components: {Reaction,Forward}
};
</script>

<template>
  
    <div class="chat-header" @click= "gotoGroupProfile">
        <img :src="'data:image/jpg;base64,' + avatar" class="user-avatar" @click="showphoto" />
        <div class="chat-info">
        <span class="chat-name">
            {{ Number(this.groupId) === 0 ? userToSend : groupname }}
        </span>
        <template v-if="Number(this.groupId) !== 0">
            <span class="group-members">
                {{ groupMembers.map(member => member.username === owner ? "io" : member.username).join(", ") }}
            </span>
        </template>
        
        <button v-if="Number(groupId) !== 0" class="add-button" @click.stop="opendModal">
            Add
          </button>
        </div>
        </div>
        
  


    <!-- Modal per la selezione utenti -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal-content">
        <h3>Search and Add Users</h3>
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Search users..."
          class="form-control mb-3"
          @input="searchUsers"
        />
        <ul class="user-list">
          <li
            v-for="user in searchResults.filter(u => u.username !== owner)"
            :key="user.userID"
            class="user-item"
            :class="{'selected-user': usertoadd.some(u => u.userID === user.userID)}"
            @click="toggleUserSelection(user)"
          >
            {{ user.username }}
        </li>
      </ul>
        <button class="confirm-button" @click="addToGroup">Add to Group</button>
        <button class="cancel-button" @click="toggleModal">Cancel</button>
      </div>
    </div>
    
    <div class="chat-messages">
      
  <div
    
    v-for="(msg, index) in messages"
    :key="index"
    :class="msg.user?.username === owner ? 'message-out' : 'message-in'"
    @click="toggleDropdown(index)"
  >
  <!-- Mostra il nome dell'utente solo se il messaggio è ricevuto e la chat è di gruppo -->
  <template v-if="Number(groupId) !== 0 && msg.user?.username !== owner">
      <div class="message-sender">{{ msg.user?.username }}</div>
    </template>
  <div v-if="dropdownIndex === index" class="dropdown-menu">
      <button v-if="msg.user?.username === owner" @click="deleteMessage(msg)">🗑️ Cancella</button>
      <button @click="replyMessage(msg)">↩️ Rispondi</button>
      <button v-if="msg.user?.username === owner" @click="showModalForward(msg)">📨 Inoltra</button>
      <button v-if="msg.user?.username !== owner" @click="showModalReaction(msg)">💬 Commenta</button>
      <button v-if="msg.user?.username !== owner" @click="showModalForward(msg)">📨 Inoltra</button>
      <button v-if="msg.message.photo!=''" @click="showphotomessageB(msg)" >Vedi FOTO</button>

    </div>
    <div class= "forwarded" v-if="msg?.message?.forwarded">Forwarded</div>
      <img class="messagephoto" v-if="msg.message.photo!=''" :src="`data:image/jpeg;base64,${msg.message.photo}`" />
    
  
      <div class="messagereply" v-if="messages.find(mex => mex.message.messageID===msg.message.linkmessage)">
        <div>
  
          <span>{{messages.find(mex => mex.message.messageID===msg.message.linkmessage).user.username}}:</span>
        </div>
        <span>{{messages.find(mex => mex.message.messageID===msg.message.linkmessage).message.txt}}</span>
      </div>
    
      
      
     
    {{ msg.message.txt }}
    
     <span v-if="msg.message.time" class="message-time">{{ msg.message.time }}</span> 
     <span v-if="msg.user?.username === owner && !msg.message.checkmark" class="checkmark">✔</span>
<span v-else-if="msg.user?.username === owner && msg.message.checkmark" class="checkmark">✔✔</span>
    <div v-if="msg.comment && msg.comment.length">
      <div v-for="comment in msg.comment" :key="comment.username">
  {{ comment.commentTXT }} - {{ comment.username === owner ? 'Tu' : comment.username }}

     

    
     
  </div>
</div>
  </div>
  </div>
<div class:="messagereply" v-if="replyId" >
  <div>
    <span>{{messages.find(mex => mex.message.messageID === replyId).user.username}}</span>
    </div>
    <div>
    <span>{{messages.find(mex => mex.message.messageID === replyId).message.txt}}</span>
  </div>
  </div> 


    <div class="chat-input">
      <input
        v-model="newMessage"
        type="text"
        placeholder="Type a message"
        @keyup.enter="check"
      />
      <input type="file" ref="fileInput" @change="sendPhoto" style="display: none;" />
      <button class="photo-button" @click="triggerFileInput">
         📷
      </button>
      <button @click="check">Send</button>
    </div>

    <div class="bigphoto" v-if="showBigPhoto">
    <img :src="`data:image/jpg;base64,${avatar}`" alt="Profile Picture" class="user-bigphoto" @click="showphoto"/>
  </div>
  <div class="bigphoto" v-if="showphotomessage">
    <img :src="`data:image/jpg;base64,${selectedMessageph}`" alt="Profile Picture" class="user-bigphoto" @click="boolshowph"/>
  </div>

    <Reaction :show="modalReaction" @close="showModalReaction" :msg="selectedMessage" :comments="comments"> >
      <template v-slot:header>
        <h3>Choose an emoji</h3>
      </template>
    </Reaction>

    <Forward :show="modalForward" @close="showModalForward" :msg="selectedMessage" :conversationlist="convList">
      <template v-slot:header>
        <h4>Forcio</h4> 
      </template>
    </Forward>
  
</template>

<style>
.messagereply {
  background-color: #f1f1f1; /* Sfondo leggermente grigio */
  border-left: 3px solid #d3d3d3; /* Bordo a sinistra per evidenziare */
  padding: 5px 10px; /* Spaziatura interna */
  margin: 5px 0; /* Spaziatura esterna */
  border-radius: 5px; /* Angoli arrotondati */
  font-size: 12px; /* Testo più piccolo */
  color: #6c757d; /* Grigio chiaro */
}

.messagereply span {
  display: block; /* Ogni elemento su una nuova riga */
  margin-bottom: 2px; /* Spaziatura tra username e messaggio */
}

.messagereply span:first-child {
  font-weight: bold; /* Username in grassetto */
  color: #495057; /* Grigio leggermente più scuro per l'username */
}
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
.messagephoto {
  width: 100px; /* Imposta una larghezza fissa per testare */
  height: auto; /* Mantieni le proporzioni */
  display: block; /* Assicurati che l'immagine sia visibile */
}
.photo-button {
  margin-left: 10px;
  padding: 10px;
  background: #f0f0f0;
  border: 1px solid #ddd;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
}

.photo-button:hover {
  background: #e0e0e0;
}
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
.forwarded {
    font-size: 12px;
    color: #667781; /* Grigio simile a WhatsApp */
    display: flex;
    align-items: center;
    gap: 4px;
    font-family: Arial, sans-serif;
    margin-bottom: 2px;
}

.forwarded::before {
    content: "\27A2"; /* Freccia simile a WhatsApp */
    font-size: 12px;
    color: #667781;
}
.chat-header {
  display: flex;
  align-items: center;
  padding: 10px;
  background: #5ed2f6;
  color: white;
}

.user-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin-right: 10px;
}

.message-time {
  font-size: 10px;
  color: #667781; /* Grigio chiaro */
  position: absolute;
  bottom: 1px;
  right: 6px;
  white-space: nowrap;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.message-in {
  position: relative;
  align-self: flex-start;
  margin: 5px 0;
  background: #e9ecef;
  padding: 10px;
  border-radius: 10px;
  max-width: 70%;
  word-wrap: break-word;
}

.message-out {
  position: relative;
  align-self: flex-end;
  margin: 5px 0;
  background: #022171;
  color: white;
  padding: 10px;
  border-radius: 10px;
  max-width: 70%;
  word-wrap: break-word;
}

.chat-input {
  display: flex;
  padding: 10px;
  background: #f1f1f1;
  border-top: 1px solid #ddd;
}

.chat-input input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
}

.chat-input button {
  margin-left: 10px;
  padding: 10px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
}
.chat-name {
  font-size: 20px; /* Cambia la dimensione del testo */
  font-weight: bold; /* Rendi il testo in grassetto */
  color: #000000; /* Colore del testo */
  margin-left: 30px; /* Distanza dall'avatar */
}
.message-sender {
  font-weight: bold;
  font-size: 14px;
  color: #555;
  margin-bottom: 3px;
}
.chat-info {
  display: flex;
  flex-direction: column;
}

.group-members {
  font-size: 12px;  /* Più piccolo del nome del gruppo */
  font-weight: bold;  /* Grassetto */
  color: #000000;  /* Colore chiaro */
  margin-top: 5px;  /* Distanza dal nome del gruppo */
  margin-left: 30px;
}

.add-button {
  margin-left: 70vh; /* Spinge il bottone a destra */
  background: #043eed;
  color: white;
  border: none;
  padding: 8px 15px;
  border-radius: 5px;
  cursor: pointer;
  font-size: 14px;
}

.add-button:hover {
  background: #218838;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  background: white;
  padding: 20px;
  border-radius: 10px;
  width: 300px;
  text-align: center;
}

.search-input {
  width: 90%;
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
}

.user-list {
  max-height: 200px;
  overflow-y: auto;
}

.user-item {
  padding: 8px;
  border: 1px solid #ddd;
  margin-bottom: 5px;
  cursor: pointer;
  border-radius: 5px;
  transition: 0.2s;
}


.selected-user {
  background-color: #007bff; /* Highlight color */
  color: white; /* Text color for contrast */
  font-weight: bold; /* Optional: make the text bold */
}
.selected-user:hover {
  background-color: #0056b3; /* Darker shade on hover */
}

.selected-users {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  margin-bottom: 10px;
}

.selected-user {
  background-color: #007bff;
  color: white;
  padding: 5px 10px;
  border-radius: 15px;
  display: flex;
  align-items: center;
}
.confirm-button {
  background: #28a745;
  color: white;
  border: none;
  padding: 10px;
  border-radius: 5px;
  cursor: pointer;
  margin-right: 5px;
}

.cancel-button {
  background: #dc3545;
  color: white;
  border: none;
  padding: 10px;
  border-radius: 5px;
  cursor: pointer;

}
.dropdown-menu {
  display: block;
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 120px;
  background: white;
  border: 2px solid red;
  
  box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.2);
  border-radius: 5px;
  z-index: 1000;
  
 
 
}
.dropdown-menu button {
  background: none;
  border: none;
  padding: 5px;
  text-align: left;
  width: 100%;
  cursor: pointer;
}

.dropdown-menu button:hover {
  background: #f0f0f0;
}

.message-in .dropdown-menu,
.message-out .dropdown-menu {
  right: 0;
}

.message-out .dropdown-menu {
  left: auto;
  right: 0;
}

.message-in .dropdown-menu {
  left: 0;
  right: auto;
}


</style>