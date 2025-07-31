# Mantel - REST API for 'a' social media.

This project is a RESTful API backend for a generic API written entirely in Go. It provides everything that a social media needs,
such as user authentication, posts, likes, follows, friendships and more.

## Work In Progress

> This project is under active development. Some features are incomplete or not yet implemented.

### Features checklist

| Feature | Status | Notes |
|---------|--------|-------|
| User Authentication | ✅ Completed | JWT-based login and registration |
| Posts | 🟡 In Progress | Basic creation works; delete and search by ID works | 
| Likes | ❌ Not started | Feature planned, not yet implemented; posts need to work properly first |
| Follows | 🧪 Implemented | Users can follow/unfollow each other; not properly tested |
| Friendships | 🧪 Implemented | Friendship relations can be formed; not properly tested |
### Legend
- ✅ Completed — Fully functional and extensively tested.
- 🟡 In Progress — Currently being worked on, not yet fully implemented.
- 🧪 Implemented — Feature works as intended, but not fully tested.
- ❌ Not started — Not implemented yet.

---

![Project Status](https://img.shields.io/badge/status-WIP-yellow)
![Go Version](https://img.shields.io/badge/go-1.20+-00ADD8?logo=go)
![Tests](https://img.shields.io/badge/tests-passing-brightgreen?logo=go)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-blue)
