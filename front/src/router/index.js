import Vue from 'vue'
import VueRouter from 'vue-router'
import axios from 'axios'

import Top from '../components/Top'
import Login from '../components/Login'
import Regist from '../components/resist'
import Questions from '../components/Questions'
import Question from '../components/Question'
import QuestionAdd from '../components/QuestionAdd'
import QuestionEdit from '../components/QuestionEdit'
import User from '../components/User'
import Users from '../components/Users'
import Admin from '../components/Admin'
import AdminQuestions from '../components/AdminQuestions'
import Logout from '../components/Logout'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Top',
    component: Top
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/regist',
    name: 'Regist',
    component: Regist
  },
  {
    path: '/question',
    name: 'Questions',
    component: Questions,
    meta: { requiresAuth: true }
  },
  {
    path: '/question/:id',
    name: 'Question',
    component: Question,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/question',
    name: 'AdminQuestions',
    component: AdminQuestions,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/question/add',
    name: 'QuestionAdd',
    component: QuestionAdd,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/question/:id/edit',
    name: 'QuestionEdit',
    component: QuestionEdit,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/user',
    name: 'Users',
    component: Users,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/user/:id',
    name: 'User',
    component: User,
    meta: { requiresAuth: true }
  },
  {
    path: '/logout',
    name: 'Logout',
    component: Logout
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    axios.post('/apis/token', { withCredentials: true }).then(() => {
      // if (res.data.token) this.$cookies.set('token', res.data.token)
      next()
    }).catch(() => {
      //next({ path: '/login', query: { redirect: to.fullPath } })
      next()
    })
  } else next()
})

export default router
