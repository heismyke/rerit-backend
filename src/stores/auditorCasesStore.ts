import { ref } from 'vue'
import { api } from '@/services/api'

export type AuditCase = {
  id: string
  property: string
  owner: string
  auditor: string
  priority: 'Low' | 'Medium' | 'High' | 'Critical'
  status: 'In Progress' | 'Pending' | 'Completed' | 'Flagged'
  started: string
  due: string
  resultStatus: 'Compliant' | 'Non-Compliant' | null
  resultNotes: string
  resultSentAt: string | null
}

export type FlaggedCase = {
  id: string
  filingId: string
  property: string
  taxpayer: string
  reason: string
  receivedAt: string
  status: 'Pending Review' | 'In Progress' | 'Resolved'
  priority: 'Low' | 'Medium' | 'High' | 'Critical'
  resultStatus: 'Compliant' | 'Non-Compliant' | null
  resultNotes: string
  resultSentAt: string | null
}

export type SuccessfulFiling = {
  id: string
  filingId: string
  property: string
  taxpayer: string
  validatedAt: string
  status: 'Validated' | 'Reviewed'
}

const auditCases = ref<AuditCase[]>([
  { id: 'AUD-2024-001', property: 'Plot 42, Victoria Island', owner: 'Emeka Okonkwo', auditor: 'John Smith', priority: 'High', status: 'In Progress', started: '2024-01-10', due: '2024-01-25', resultStatus: null, resultNotes: '', resultSentAt: null },
  { id: 'AUD-2024-002', property: 'Block 7, Lekki Phase 2', owner: 'Adaobi Nnamdi', auditor: 'Sarah Johnson', priority: 'Medium', status: 'Pending', started: '2024-01-12', due: '2024-01-30', resultStatus: null, resultNotes: '', resultSentAt: null },
])

const flaggedCases = ref<FlaggedCase[]>([
  { id: 'FLG-2024-011', filingId: 'FCT-IRS-00211', property: 'Plot 18, Wuse II', taxpayer: 'Nwosu Holdings', reason: 'Declared rent below benchmark', receivedAt: '2024-01-18', status: 'Pending Review', priority: 'High', resultStatus: null, resultNotes: '', resultSentAt: null },
  { id: 'FLG-2024-012', filingId: 'FCT-IRS-00225', property: 'Block 4, Maitama', taxpayer: 'Ayo Martins', reason: 'Declared rent below benchmark', receivedAt: '2024-01-18', status: 'Pending Review', priority: 'Medium', resultStatus: null, resultNotes: '', resultSentAt: null },
  { id: 'FLG-2024-013', filingId: 'FCT-IRS-00231', property: 'Unit 12, Gwarinpa', taxpayer: 'Saka Ventures', reason: 'Declared rent below benchmark', receivedAt: '2024-01-19', status: 'In Progress', priority: 'High', resultStatus: null, resultNotes: '', resultSentAt: null },
  { id: 'FLG-2024-014', filingId: 'FCT-IRS-00237', property: 'Plot 9, Asokoro', taxpayer: 'Dara Okafor', reason: 'Declared rent below benchmark', receivedAt: '2024-01-20', status: 'In Progress', priority: 'Critical', resultStatus: 'Non-Compliant', resultNotes: 'Under-declaration confirmed. Notice issued.', resultSentAt: '2024-01-21' },
  { id: 'FLG-2024-004', filingId: 'FCT-IRS-00004', property: 'Plot 8, Banana Island', taxpayer: 'Folake Adeyemi', reason: 'Declared rent below benchmark', receivedAt: '2024-01-18', status: 'In Progress', priority: 'Critical', resultStatus: 'Non-Compliant', resultNotes: 'Under-declared rental income. Notice of non-compliance issued.', resultSentAt: '2024-01-18' },
])

const successfulFilings = ref<SuccessfulFiling[]>([
  { id: 'SUC-2024-101', filingId: 'FCT-IRS-00311', property: 'Plot 12, Jabi', taxpayer: 'Kola Ibrahim', validatedAt: '2024-01-18', status: 'Validated' },
  { id: 'SUC-2024-102', filingId: 'FCT-IRS-00318', property: 'Block 5, Garki', taxpayer: 'Laila Musa', validatedAt: '2024-01-18', status: 'Validated' },
  { id: 'SUC-2024-103', filingId: 'FCT-IRS-00327', property: 'Unit 3, Wuse I', taxpayer: 'Prime Estates Ltd', validatedAt: '2024-01-19', status: 'Validated' },
  { id: 'SUC-2024-104', filingId: 'FCT-IRS-00333', property: 'Plot 6, Maitama', taxpayer: 'Uche Nwankwo', validatedAt: '2024-01-20', status: 'Validated' },
  { id: 'SUC-2024-003', filingId: 'FCT-IRS-00003', property: '15 Admiralty Way, Lekki', taxpayer: 'Chidi Okafor', validatedAt: '2024-01-19', status: 'Validated' },
])

type AuditorSnapshot = {
  auditCases: AuditCase[]
  flaggedCases: FlaggedCase[]
  successfulFilings: SuccessfulFiling[]
}

const isLoaded = ref(false)

const loadAuditorCases = async () => {
  if (isLoaded.value) return
  const [remoteAuditCases, remoteFlaggedCases, remoteSuccessfulFilings] = await Promise.all([
    api.getAuditCases<AuditCase[]>(),
    api.getFlaggedCases<FlaggedCase[]>(),
    api.getSuccessfulFilings<SuccessfulFiling[]>(),
  ])
  auditCases.value = remoteAuditCases
  flaggedCases.value = remoteFlaggedCases
  successfulFilings.value = remoteSuccessfulFilings
  isLoaded.value = true
}

const buildFilingId = (id: string) => `FCT-IRS-${id.replace(/\D/g, '').padStart(5, '0')}`

const moveAuditToSuccessful = (c: AuditCase) => {
  successfulFilings.value.unshift({
    id: `SUC-${c.id.split('-').slice(1).join('-')}`,
    filingId: buildFilingId(c.id),
    property: c.property,
    taxpayer: c.owner,
    validatedAt: new Date().toISOString().split('T')[0],
    status: 'Validated',
  })
}

const moveAuditToFlagged = (c: AuditCase) => {
  flaggedCases.value.unshift({
    id: `FLG-${c.id.split('-').slice(1).join('-')}`,
    filingId: buildFilingId(c.id),
    property: c.property,
    taxpayer: c.owner,
    reason: 'Declared rent below benchmark',
    receivedAt: new Date().toISOString().split('T')[0],
    status: 'Pending Review',
    priority: c.priority,
    resultStatus: c.resultStatus,
    resultNotes: c.resultNotes,
    resultSentAt: c.resultSentAt,
  })
}

const createAuditCase = async (draft: Partial<AuditCase>) => {
  const created = await api.createAuditCase<AuditCase>(draft)
  auditCases.value.unshift(created)
  return created
}

const updateAuditCase = async (id: string, patch: AuditCase) => {
  const updated = await api.updateAuditCase<AuditCase>(id, patch)
  const index = auditCases.value.findIndex(c => c.id === id)
  if (index !== -1) auditCases.value[index] = updated
  return updated
}

const deleteAuditCase = async (id: string) => {
  await api.deleteAuditCase(id)
  auditCases.value = auditCases.value.filter(c => c.id !== id)
}

const sendAuditResult = async (id: string, status: NonNullable<AuditCase['resultStatus']>, notes: string) => {
  const snapshot = await api.sendAuditResult<AuditorSnapshot>(id, { status, notes })
  auditCases.value = snapshot.auditCases
  flaggedCases.value = snapshot.flaggedCases
  successfulFilings.value = snapshot.successfulFilings
}

const updateFlaggedCase = async (id: string, patch: FlaggedCase) => {
  const updated = await api.updateFlaggedCase<FlaggedCase>(id, patch)
  const index = flaggedCases.value.findIndex(f => f.id === id)
  if (index !== -1) flaggedCases.value[index] = updated
  return updated
}

const sendFlaggedResult = async (id: string, status: NonNullable<FlaggedCase['resultStatus']>, notes: string) => {
  const updated = await api.sendFlaggedResult<FlaggedCase>(id, { status, notes })
  const index = flaggedCases.value.findIndex(f => f.id === id)
  if (index !== -1) flaggedCases.value[index] = updated
  return updated
}

const updateSuccessfulFiling = async (id: string, patch: SuccessfulFiling) => {
  const updated = await api.updateSuccessfulFiling<SuccessfulFiling>(id, patch)
  const index = successfulFilings.value.findIndex(f => f.id === id)
  if (index !== -1) successfulFilings.value[index] = updated
  return updated
}

export const useAuditorCasesStore = () => {
  return {
    auditCases,
    flaggedCases,
    successfulFilings,
    loadAuditorCases,
    moveAuditToSuccessful,
    moveAuditToFlagged,
    createAuditCase,
    updateAuditCase,
    deleteAuditCase,
    sendAuditResult,
    updateFlaggedCase,
    sendFlaggedResult,
    updateSuccessfulFiling,
  }
}
