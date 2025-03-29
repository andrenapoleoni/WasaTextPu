<script>
export default {
  props: {
    conversationlist: Array,
    show: Boolean,
    msg: Object,
  },
  data() {
    return {
      searchPrivateTerm: "",
      searchGroupTerm: "",
      searchPrivateResults: [],
      searchGroupResults: [],
      errormsg: "",
      owner: sessionStorage.username,
      selectedUsers: [],
      selectedGroups: [],
      activeTab: 'private', // Gestisce quale tab è attivo
      convId: parseInt(this.$route.params.conversation),
    };
  },
  methods: {

    async ForwardMessage(){
      if (this.selectedUsers.length === 0 && this.selectedGroups.length === 0) {
        this.errormsg = "Seleziona almeno un utente o un gruppo per inoltrare il messaggio.";
        console.warn(this.errormsg);
        return; // Esce dalla funzione senza eseguire il try
      }
      console.log("convid:", this.convId);
      console.log("Forwarding message:", this.msg.messageID);
      console.log("Selected users:", this.selectedUsers);
      console.log("Selected groups:", this.selectedGroups);
      try {
        console.log("dentro try");
        for (const conv of this.selectedUsers) {  // Cambiato da "in" a "of"
          console.log("dentro for 1");
          console.log("username:", conv.username);
          let response=await this.$axios.post(`/user/${sessionStorage.userID}/conversation/${this.convId}/messages/${this.msg.messageID}`, {},
          {
            params: {dest: conv.username },  // Adesso conv è l'oggetto corretto
            headers: { Authorization: sessionStorage.token }
          });
        }

        for (const conv of this.selectedGroups) {  // Cambiato da "in" a "of"
          console.log("dentro for 2");
          console.log("groupID:", conv.group.groupID);
          console.log("group:", conv.groupID);
          let response =await this.$axios.post(`/user/${sessionStorage.userID}/conversation/${this.convId}/messages/${this.msg.messageID}`, {},
           {
            params:{ dest: conv.group.groupID },  // Adesso conv è l'oggetto corretto
            headers: { Authorization: sessionStorage.token },
          });


        }
        this.closeModal();
  } catch (e) {
    this.errormsg = e.toString();
  }
    },
    async searchUsers() {
      if (!this.searchPrivateTerm) {
        this.searchPrivateResults = [];
        return;
      }
      try {
        const response = await this.$axios.get(`/user/${sessionStorage.userID}`, {
          params: { search: this.searchPrivateTerm },
          headers: { Authorization: sessionStorage.token },
        });
        this.searchPrivateResults = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
    },
    filterGroups() {
      console.log("conversationlist:", this.conversationlist);
      if (!Array.isArray(this.conversationlist)) {
        this.errormsg = "conversationlist is not an array";
        return;
      }
      this.searchGroupResults = this.conversationlist.filter(convo => convo.group.groupID !== 0);
      console.log("Groups:", this.searchGroupResults);
    },
    selectUser(user) {
      const index = this.selectedUsers.findIndex(u => u.username === user.username);
      if (index === -1) {
        this.selectedUsers.push(user);
      } else {
        this.selectedUsers.splice(index, 1);
      }
    },
    selectGroup(group) {
      console.log("Selezionato:", group.group.groupname, "ID:", group.group.groupID);
      
      // Evidenzia il gruppo selezionato
      if (this.selectedGroupID === group.group.groupID) {
        this.selectedGroupID = null;
      } else {
        this.selectedGroupID = group.group.groupID;
      }

      // Gestisci l'aggiunta o rimozione dal gruppo selezionato
      const index = this.selectedGroups.findIndex((g) => g.group.groupID === group.group.groupID);
      if (index === -1) {
        this.selectedGroups.push(group);
      } else {
        this.selectedGroups.splice(index, 1);
      }
  },
    
    closeModal() {
      this.searchPrivateTerm = "";
      this.searchPrivateResults = [];
      this.selectedUsers = [];
      this.selectedGroups = [];
      window.location.reload();
      this.$emit("close");
    },
  },
};
</script>

<template>
  <div v-if="show" class="modal-overlay">
    <div class="modal-content">
      <h3>Forward</h3>
      <div class="tab-buttons">
        <button class="btn" :class="{ active: activeTab === 'private' }" @click="activeTab = 'private'">Private</button>
        <button class="btn" :class="{ active: activeTab === 'group' }" @click="activeTab = 'group'; filterGroups()">Group</button>
      </div>

      <!-- Sezione utenti e gruppi selezionati -->
      <div v-if="selectedUsers.length || selectedGroups.length" class="selected-container">
        <div v-if="selectedUsers.length" class="selected-users">
          <h4>Selected Users:</h4>
          <span v-for="user in selectedUsers" :key="user.userID" class="selected-item">
            {{ user.username }}
            <span class="remove-item" @click="selectUser(user)">✖</span>
          </span>
        </div>
        <div v-if="selectedGroups.length" class="selected-groups">
          <h4>Selected Groups:</h4>
          <span v-for="group in selectedGroups" :key="group.group.groupID" class="selected-item">
            {{ group.group.groupname }}
            <span class="remove-item" @click="selectGroup(group)">✖</span>
          </span>
        </div>
      </div>

      <!-- Sezione Private -->
      <div v-if="activeTab === 'private'">
        <input
          type="text"
          v-model="searchPrivateTerm"
          placeholder="Search for users"
          class="form-control mb-3"
          @input="searchUsers"
        />
        <ul class="list-group">
          <li
            v-for="user in searchPrivateResults.filter(user => user.username !== owner)"
            :key="user.userID"
            class="list-group-item"
            :class="{ 'selected-item': selectedUsers.some(u => u.userID === user.userID) }"
            @click="selectUser(user)"
          >
            {{ user.username }}
          </li>
        </ul>
      </div>

      <!-- Sezione Group -->
      <div v-if="activeTab === 'group'">
        <ul class="list-group">
          <li v-for="group in searchGroupResults" :key="group.group.groupID" @click="selectGroup(group)">
            {{ group.group.groupname }}
          </li>
        </ul>
      </div>
      
      <div button class="btn btn-primary" @click="closeModal">Close</div>
      <div button class="btn btn-primary" @click="ForwardMessage">Send</div>
    </div>
  </div>
</template>

<style>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}
.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 400px;
  text-align: center;
}
.tab-buttons {
  display: flex;
  justify-content: space-around;
  margin-bottom: 10px;
}
.btn {
  padding: 8px 16px;
  cursor: pointer;
  border: none;
  background: #ccc;
}
.btn.active {
  background: #007bff;
  color: white;
}
.selected-container {
  margin-bottom: 10px;
}
.selected-item {
  display: inline-block;
  background: #007bff;
  color: white;
  padding: 5px 10px;
  margin: 5px;
  border-radius: 5px;
}
.remove-item {
  cursor: pointer;
  margin-left: 5px;
}
</style>
