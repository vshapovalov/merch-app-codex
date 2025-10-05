import { createApp } from 'vue';
import PrimeVue from 'primevue/config';
import ToastService from 'primevue/toastservice';
import Button from 'primevue/button';
import Menubar from 'primevue/menubar';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Dropdown from 'primevue/dropdown';
import Calendar from 'primevue/calendar';
import Password from 'primevue/password';
import Card from 'primevue/card';
import Toast from 'primevue/toast';
import Toolbar from 'primevue/toolbar';
import Divider from 'primevue/divider';
import InputNumber from 'primevue/inputnumber';

import '@primevue/themes/lara';
import 'primeicons/primeicons.css';
import 'primeflex/primeflex.css';
import './style.css';

import App from './App.vue';
import router from './router';
import { createPinia } from 'pinia';

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(PrimeVue, { ripple: true });
app.use(ToastService);

app.component('Button', Button);
app.component('Menubar', Menubar);
app.component('DataTable', DataTable);
app.component('Column', Column);
app.component('Dialog', Dialog);
app.component('InputText', InputText);
app.component('Textarea', Textarea);
app.component('Dropdown', Dropdown);
app.component('Calendar', Calendar);
app.component('Password', Password);
app.component('Card', Card);
app.component('Toast', Toast);
app.component('Toolbar', Toolbar);
app.component('Divider', Divider);
app.component('InputNumber', InputNumber);

app.mount('#app');
