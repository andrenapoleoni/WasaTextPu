<script>

import Group from "../components/GroupModal.vue";
import Private from "../components/PrivateModal.vue";

export default {
  data: function() {
    return {
      owner: sessionStorage.username, 
      ownerID: sessionStorage.userID,   
      errormsg: null,
      some_data: [],
      searchTerm: "",
      searchResults: [],
      existingChatWarning: "",
      showGroupChatModal: false,
      groupChatMembers: [],
      showPrivateChatModal: false,
      showDropdown: false,
      intervalId: null,
      

    };
  },
  methods: {
    async getMyconversations() {
      this.errormsg = null;
      console.log("USer:", sessionStorage.userID,"token:", sessionStorage.token);
      try {
        let response = await this.$axios.get(`/user/${sessionStorage.userID}/conversation`, {
          headers: { 'Authorization': sessionStorage.token }
      });

      this.some_data = response.data.map(chat => {
      return {
        ...chat,
        isGroup: !!chat.group, // Se `group` esiste, allora è una chat di gruppo
      };
      });
      console.log("My conversations:", this.some_data);
     
    

      } catch (e) {
        this.errormsg = e.toString();
  }
},

    handleConversationClick(response) {
      localStorage.clear();
      
    if (response.group.groupID!==0) {
      console.log("Group chat:", response.group.groupname);
      localStorage.groupname = response.group.groupname;
      localStorage.groupId = response.group.groupID;
      localStorage.groupMembers = JSON.stringify(response.groupUsers);
      localStorage.photo = response.group.photo;
    }else {
      console.log("Private chat:", response.user.username);
      localStorage.groupId = 0;
      localStorage.username = response.user.username;
      localStorage.photo = response.user.photo;
    }
    localStorage.conversationList = JSON.stringify(this.some_data);
    console.log("Conversation list:", localStorage.conversationList);
    this.$router.push(`/conversation/${response.conversation.conversationID}`);
  },
  
   
    handleCreatePrivateChatModalToggle() {
      console.log("Create private chat",this.some_data);
      this.showPrivateChatModal = !this.showPrivateChatModal;
    },
    handleCreateGroupModalToggle() {
      this.showGroupChatModal = !this.showGroupChatModal;
   
    },
    toggleDropdown(){
      this.showDropdown = !this.showDropdown;
    }
},
    

    

  mounted() {
    if (!sessionStorage.token) {
      this.$router.push("/session");
      return;
    }
    this.getMyconversations();
    document.addEventListener("click", this.closeDropdown);
    this.intervalId = setInterval(async () => {
      clearInterval(this.intervalId);
      await this.getMyconversations();
      this.intervalId = setInterval(this.getMyconversations, 7000);
    }, 7000);
  },
  beforeUnmount() {
    document.removeEventListener("click", this.closeDropdown);
  },
  components: {Group,Private},

};


</script>

<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Home page</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="mounted">Refresh</button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">Export</button>
        </div>
        <div class="btn-group me-2 position-relative" @click.stop>
          <button type="button" class="btn btn-sm btn-outline-primary" @click="toggleDropdown">New Chat</button>
          <div v-if="showDropdown" class="dropd-menu show dropd-menu-end">
            <button class="dropd-item" @click="handleCreateGroupModalToggle">New Group Chat</button>
            <button class="dropd-item" @click="handleCreatePrivateChatModalToggle">New Private Chat</button>
          </div>
        </div>
      </div>
    </div>


    <!-- Error message -->
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <!-- Warning for existing chat -->
    <div v-if="existingChatWarning" class="alert alert-warning">
      {{ existingChatWarning }}
    </div>

    
    <div 
  class="conversations" 
  v-for="response in some_data" 
  :key="response.conversation.conversationID"
  @click="handleConversationClick(response)"
  style="cursor: pointer;"
>
  <!-- Se esiste un gruppo, mostra il nome del gruppo -->
   <div v-if="response.group && response.group.groupname" class="conversation-item">
  <img 
      v-if="response.group.photo" 
      :src="'data:image/jpeg;base64,' + response.group.photo" 
      alt="Group Photo" 
      class="conversation-photo"
    />
  <span  class="conversation-name">{{response.group.groupname }}</span>
  <span v-if="response.lastmessage.mexTXT" class="conversation-lastmessage"> Last Message: {{response.lastmessage.mexTXT }}</span>
  <span v-if="response.lastmessage.mexUser" class="lastmex-user">By: {{ response.lastmessage.mexUser===owner ? 'Tu': response.lastmessage.mexUser   }}</span>
  <span v-if="response.lastmessage.mexTime" class="lastmex-user">at: {{ response.lastmessage.mexTime}}</span>
</div>
  
  <!-- Altrimenti, mostra il nome dell'utente -->
  <div v-else class="conversation-item">
  <img 
      v-if="response.user.photo" 
      :src="'data:image/jpeg;base64,' + response.user.photo" 
      alt="User Photo" 
      class="conversation-photo"
    />
  <span v-if="response.user.username" class="conversation-name">{{response.user.username }}</span>
  <span v-if="response.lastmessage.mexTXT" class=conversation-lastmessage>  Last Message: {{response.lastmessage.mexTXT}}</span>
  <span v-if="response.lastmessage.mexUser" class="lastmex-user">By: {{ response.lastmessage.mexUser===owner ? 'Tu': response.lastmessage.mexUser  }}</span>
  <span v-if="response.lastmessage.mexTime" class="lastmex-user">at: {{ response.lastmessage.mexTime}}</span>
</div>
  <hr>


</div>
  </div>
  
   <!-- Modale utilizzato per la creazione di un nuovo gruppo -->
   <Group :show="showGroupChatModal" @close="handleCreateGroupModalToggle" title="search">
        <template v-slot:header>
          <h3>Select users</h3>
        </template>
      </Group>

      <!-- Modal for private chat -->
    <Private :show="showPrivateChatModal" @close="handleCreatePrivateChatModalToggle" title="search" :somedata="some_data">
    <template v-slot:header>
          <h3>Select users</h3>
        </template>
      </Private>

  


</template>

<style>
.position-relative {
  position: relative;
}
.conversation-item{
  display: flex;
  align-items: center;
  padding: 10px;
  gap: 10px;
}
.conversation-name {
  font-weight: bold;
  font-size: 16px;
  color: #333;
}
.lastmex-user {
  font-size: 14px;
  color: #666;
  margin-top: 5px; /* Spazio tra il nome e l'ultimo messaggio */
  margin-left: 100px;
}
.conversation-lastmessage {
  font-size: 14px;
  color: #666;
  margin-top: 5px; /* Spazio tra il nome e l'ultimo messaggio */
  margin-left: 100px;
}
.conversation-photo{
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 10px;
}
.dropd-menu {
  position: absolute;
  top: 100%;
  right: 0; 
  left: auto;
  z-index: 1000;
  display: block;
  float: left;
  min-width: 10rem;
  padding: .5rem 0;
  margin: .125rem 0 0;
  font-size: 1rem;
  color: #212529;
  text-align: left;
  background-color: #fff;
  border: 1px solid rgba(0, 0, 0, .15);
  border-radius: .25rem;
  box-shadow: 0 .5rem 1rem rgba(0, 0, 0, .175);
}

.dropd-item {
  display: block;
  width: 100%;
  padding: .25rem 1.5rem;
  clear: both;
  font-weight: 400;
  color: #212529;
  text-align: inherit;
  white-space: nowrap;
  background-color: transparent;
  border: 0;
  cursor: pointer;
}

.dropd-item:hover {
  color: #fff;
  background-color: #ddeb1c;
}

/* Modal Styles */
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
  z-index: 1050;
}

.modal-content {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 400px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  text-align: center;
}

.modal-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

/* List of matching users */
.list-group {
  max-height: 200px;
  overflow-y: auto;
  margin-bottom: 10px;
}

.list-group-item {
  cursor: pointer;
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
  background-color: #1973d3;
  color: white;
  padding: 5px 10px;
  border-radius: 15px;
  display: flex;
  align-items: center;
}

.remove-user {
  margin-left: 8px;
  cursor: pointer;
  font-weight: bold;
}

</style>