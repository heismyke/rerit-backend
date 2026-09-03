# FCT-IRS - Real Estate Revenue & Information System

<p align="center">
  <img src="public/fct-irs.png" alt="FCT-IRS Logo" width="80" />
</p>

<p align="center">
  <strong>FCT Internal Revenue Service</strong><br>
  Real Estate Tax & Revenue Collection Platform
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js" alt="Vue">
  <img src="https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript" alt="TypeScript">
  <img src="https://img.shields.io/badge/Vite-5-646CFF?logo=vite" alt="Vite">
  <img src="https://img.shields.io/badge/Tailwind-4-38B2AC?logo=tailwind-css" alt="Tailwind">
</p>

---

## Mission

> How can the government ensure that all properties are registered, correctly valued, assessed for tax, audited, and enforced efficiently, while maximizing compliance and minimizing revenue leakage?

FCT-IRS is designed to address this challenge by providing a comprehensive platform for real estate revenue administration in Nigeria.

---

## Features

### Multi-Role System
ReRiT provides role-specific dashboards for 5 different user types:

| Role | Description |
|------|-------------|
| **Admin** | Full system control, rules management, user management, compliance oversight |
| **Auditor** | Property verification, audit case management, land registry oversight |
| **Compliance Officer** | Investigation of flagged properties, compliance notes, enforcement actions |
| **Real Estate Developer (Tax Payer)** | Property registration, tax payments, notice viewing |
| **Surveyor** | Property surveys, submission management |

### Core Functionality
- **Property Management** - Register, verify, and track real estate properties
- **Tax Collection** - Revenue tracking and payment processing
- **Risk Assessment Engine** - Automated risk scoring (0-100) with Low/Medium/High/Critical levels
- **Case Management** - Track investigations, fraud cases, and compliance issues
- **Survey System** - Property surveys and submissions
- **Land Registry** - Centralized property ownership records
- **Reporting** - Generate and download system reports
- **Notifications** - Real-time alerts and updates

### Risk Scoring Engine
The built-in risk assessment engine evaluates properties based on:
- Property value (N500M+ = +20 points)
- Property type (Commercial = +10 points)
- Ownership history changes (3+ = +15 points)
- Value discrepancies (Flagged = +25 points)
- Survey status (Pending/Flagged = +15 points)
- Payment history (Late/Default = +20 points)

---

## Tech Stack

| Technology | Purpose |
|------------|---------|
| **Vue 3** | Frontend framework with Composition API |
| **TypeScript** | Type-safe development |
| **Vite** | Build tool and dev server |
| **Vue Router 4** | Client-side routing |
| **Tailwind CSS 4** | Utility-first styling |
| **Chart.js + vue-chartjs** | Data visualization |
| **Pinia** | State management |

---

## Getting Started

### Prerequisites
- Node.js 18+
- npm or yarn

### Installation

```bash
# Clone the repository
git clone https://github.com/heismyke/rerit-s.git
cd rerit-s

# Install dependencies
npm install

# Start development server
npm run dev
```

The app will be available at `http://localhost:5174`

### Build for Production

```bash
npm run build
```

Preview production build:
```bash
npm run preview
```

---

## Project Structure

```
src/
├── components/
│   └── Sidebar.vue          # Shared navigation sidebar
├── data/
│   └── roles.ts             # Role definitions
├── router/
│   ├── index.ts             # Router configuration
│   └── public.ts            # Public routes
├── stores/
│   └── index.ts             # State management
├── types/
│   └── index.ts             # TypeScript types
├── utils/
│   └── riskScoring.ts       # Risk assessment engine
├── views/
│   ├── Home.vue             # Landing/role selection
│   └── dashboards/
│       ├── admin/           # Admin dashboard pages
│       ├── auditor/         # Auditor dashboard pages
│       ├── compliance/      # Compliance officer pages
│       ├── developer/       # Tax payer pages
│       └── surveyor/        # Surveyor pages
├── App.vue
└── main.ts
```

---

## NRS Brand Guidelines

| Element | Value |
|---------|-------|
| Primary Color | `#B90B0B` (NRS Red) |
| Sidebar Background | `#EEEEEE` |
| Off-White Background | `#f5f7fa` |
| Font | System sans-serif |

---

## Available Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start development server |
| `npm run build` | Build for production |
| `npm run preview` | Preview production build |
| `npm run lint` | Run ESLint |
| `npm run typecheck` | Run TypeScript type checking |

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is proprietary to the National Revenue System (NRS), Nigeria.

---

## Support

For issues or feature requests, please contact the NRS development team.

# rerit-backend
