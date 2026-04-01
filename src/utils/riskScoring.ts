export const calculateRiskScore = (property: {
  declaredValue?: number
  propertyType?: string
  location?: string
  ownershipHistory?: number
  hasDiscrepancy?: boolean
  surveyStatus?: string
  paymentHistory?: string
}): { score: number; level: 'Low' | 'Medium' | 'High' | 'Critical'; factors: string[] } => {
  let score = 30
  const factors: string[] = []

  if (property.declaredValue && property.declaredValue > 500000000) {
    score += 20
    factors.push('High value property (N500M+)')
  }

  if (property.propertyType === 'Commercial') {
    score += 10
    factors.push('Commercial property')
  }

  if (property.ownershipHistory && property.ownershipHistory > 2) {
    score += 15
    factors.push('Multiple ownership changes')
  }

  if (property.hasDiscrepancy) {
    score += 25
    factors.push('Value discrepancy detected')
  }

  if (property.surveyStatus === 'Pending' || property.surveyStatus === 'Flagged') {
    score += 15
    factors.push('Survey verification needed')
  }

  if (property.paymentHistory === 'Late' || property.paymentHistory === 'Default') {
    score += 20
    factors.push('Poor payment history')
  }

  score = Math.min(100, score)

  let level: 'Low' | 'Medium' | 'High' | 'Critical'
  if (score < 30) level = 'Low'
  else if (score < 50) level = 'Medium'
  else if (score < 75) level = 'High'
  else level = 'Critical'

  return { score, level, factors }
}

export const getRiskColor = (level: string) => {
  switch (level) {
    case 'Low': return 'bg-green-50 text-green-700'
    case 'Medium': return 'bg-yellow-50 text-yellow-700'
    case 'High': return 'bg-orange-50 text-orange-700'
    case 'Critical': return 'bg-red-50 text-red-700'
    default: return 'bg-gray-100 text-gray-700'
  }
}

export const getRiskScoreColor = (score: number) => {
  if (score < 30) return 'text-green-600'
  if (score < 50) return 'text-yellow-600'
  if (score < 75) return 'text-orange-600'
  return 'text-red-600'
}
