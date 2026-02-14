export interface Company {
  id: string;
  name: string;
  address?: string;
  phone?: string;
  email?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Worksite {
  id: string;
  companyId: string;
  name: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  createdAt: string;
  updatedAt: string;
}

export interface Worker {
  id: string;
  authSubject: string;
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  createdAt: string;
  updatedAt: string;
}

export type WorkerRole = "worker" | "company_admin" | "site_admin";
export type MembershipStatus = "active" | "inactive";

export interface WorkerCompany {
  workerId: string;
  companyId: string;
  role: WorkerRole;
  status: MembershipStatus;
  joinedAt: string;
}

export interface Certificate {
  id: string;
  workerId: string;
  name: string;
  issuingBody?: string;
  certificateNumber?: string;
  issuedDate?: string;
  expiryDate?: string;
  createdAt: string;
  updatedAt: string;
}

export type ShiftStatus = "open" | "assigned" | "in_progress" | "completed" | "cancelled";

export interface Shift {
  id: string;
  worksiteId: string;
  createdBy: string;
  title: string;
  description?: string;
  startTime: string;
  endTime: string;
  status: ShiftStatus;
  createdAt: string;
  updatedAt: string;
}

export type AssignmentStatus = "offered" | "accepted" | "declined" | "completed";

export interface ShiftAssignment {
  id: string;
  shiftId: string;
  workerId: string;
  status: AssignmentStatus;
  assignedAt: string;
  respondedAt?: string;
}

export interface ShiftReportTemplate {
  id: string;
  companyId: string;
  name: string;
  fields: string;
  createdAt: string;
  updatedAt: string;
}

export interface ShiftReport {
  id: string;
  shiftId: string;
  workerId: string;
  templateId?: string;
  data: string;
  submittedAt: string;
}

export interface LocationCheckIn {
  id: string;
  workerId: string;
  shiftId?: string;
  latitude: number;
  longitude: number;
  recordedAt: string;
}

export type AlarmStatus = "raised" | "acknowledged" | "resolved";

export interface Alarm {
  id: string;
  workerId: string;
  shiftId?: string;
  latitude?: number;
  longitude?: number;
  message?: string;
  status: AlarmStatus;
  raisedAt: string;
  acknowledgedAt?: string;
  resolvedAt?: string;
}
