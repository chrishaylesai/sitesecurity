"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams } from "next/navigation";
import {
  PageHeader,
  DataTable,
  Badge,
  Button,
  Modal,
  Input,
} from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type { Worker, Certificate, WorkerCompany, Company } from "@/lib/types";

function formatDate(dateStr?: string): string {
  if (!dateStr) return "\u2014";
  return new Date(dateStr).toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

function isExpired(dateStr?: string): boolean {
  if (!dateStr) return false;
  return new Date(dateStr) < new Date();
}

function roleBadgeVariant(role: string): string {
  switch (role) {
    case "company_admin":
      return "info";
    case "site_admin":
      return "warning";
    default:
      return "default";
  }
}

function statusBadgeVariant(status: string): string {
  switch (status) {
    case "active":
      return "success";
    case "inactive":
      return "danger";
    default:
      return "default";
  }
}

function roleLabel(role: string): string {
  switch (role) {
    case "company_admin":
      return "Company Admin";
    case "site_admin":
      return "Site Admin";
    case "worker":
      return "Worker";
    default:
      return role;
  }
}

export default function WorkerDetailPage() {
  const params = useParams();
  const workerId = params.id as string;

  const [worker, setWorker] = useState<Worker | null>(null);
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [memberships, setMemberships] = useState<WorkerCompany[]>([]);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Edit worker modal
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [editForm, setEditForm] = useState({
    firstName: "",
    lastName: "",
    email: "",
    phone: "",
    authSubject: "",
  });
  const [editSubmitting, setEditSubmitting] = useState(false);

  // Add certificate modal
  const [certModalOpen, setCertModalOpen] = useState(false);
  const [certForm, setCertForm] = useState({
    name: "",
    issuingBody: "",
    certificateNumber: "",
    issuedDate: "",
    expiryDate: "",
  });
  const [certSubmitting, setCertSubmitting] = useState(false);

  const fetchWorker = useCallback(() => {
    return apiClient<Worker>(`/api/v1/workers/${workerId}`);
  }, [workerId]);

  const fetchCertificates = useCallback(() => {
    return apiClient<Certificate[]>(`/api/v1/workers/${workerId}/certificates`);
  }, [workerId]);

  const fetchMemberships = useCallback(() => {
    return apiClient<WorkerCompany[]>(`/api/v1/workers/${workerId}/memberships`);
  }, [workerId]);

  const fetchCompanies = useCallback(() => {
    return apiClient<Company[]>("/api/v1/companies");
  }, []);

  const loadAll = useCallback(() => {
    setLoading(true);
    Promise.all([fetchWorker(), fetchCertificates(), fetchMemberships(), fetchCompanies()])
      .then(([w, certs, membs, comps]) => {
        setWorker(w);
        setCertificates(certs);
        setMemberships(membs);
        setCompanies(comps);
        setEditForm({
          firstName: w.firstName,
          lastName: w.lastName,
          email: w.email,
          phone: w.phone ?? "",
          authSubject: w.authSubject,
        });
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [fetchWorker, fetchCertificates, fetchMemberships, fetchCompanies]);

  useEffect(() => {
    loadAll();
  }, [loadAll]);

  const companyLookup = new Map(companies.map((c) => [c.id, c.name]));

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setEditSubmitting(true);
    try {
      const updated = await apiClient<Worker>(`/api/v1/workers/${workerId}`, {
        method: "PUT",
        body: JSON.stringify({
          firstName: editForm.firstName,
          lastName: editForm.lastName,
          email: editForm.email,
          phone: editForm.phone || undefined,
          authSubject: editForm.authSubject,
        }),
      });
      setWorker(updated);
      setEditModalOpen(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update worker");
    } finally {
      setEditSubmitting(false);
    }
  };

  const handleCertSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setCertSubmitting(true);
    try {
      await apiClient<Certificate>(`/api/v1/workers/${workerId}/certificates`, {
        method: "POST",
        body: JSON.stringify({
          name: certForm.name,
          issuingBody: certForm.issuingBody || undefined,
          certificateNumber: certForm.certificateNumber || undefined,
          issuedDate: certForm.issuedDate || undefined,
          expiryDate: certForm.expiryDate || undefined,
        }),
      });
      setCertModalOpen(false);
      setCertForm({ name: "", issuingBody: "", certificateNumber: "", issuedDate: "", expiryDate: "" });
      const certs = await fetchCertificates();
      setCertificates(certs);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add certificate");
    } finally {
      setCertSubmitting(false);
    }
  };

  if (loading) {
    return <div className="text-sm text-gray-500">Loading worker details...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load worker: {error}
      </div>
    );
  }

  if (!worker) {
    return (
      <div className="rounded-md bg-yellow-50 p-4 text-sm text-yellow-700">
        Worker not found.
      </div>
    );
  }

  const certificateColumns = [
    { key: "name", header: "Name" },
    {
      key: "issuingBody",
      header: "Issuing Body",
      render: (c: Certificate) => c.issuingBody ?? "\u2014",
    },
    {
      key: "certificateNumber",
      header: "Certificate Number",
      render: (c: Certificate) => c.certificateNumber ?? "\u2014",
    },
    {
      key: "issuedDate",
      header: "Issued Date",
      render: (c: Certificate) => formatDate(c.issuedDate),
    },
    {
      key: "expiryDate",
      header: "Expiry Date",
      render: (c: Certificate) => (
        <span className={isExpired(c.expiryDate) ? "font-medium text-red-600" : ""}>
          {formatDate(c.expiryDate)}
        </span>
      ),
    },
  ];

  const membershipColumns = [
    {
      key: "companyName",
      header: "Company Name",
      render: (m: WorkerCompany) => companyLookup.get(m.companyId) ?? m.companyId,
    },
    {
      key: "role",
      header: "Role",
      render: (m: WorkerCompany) => (
        <Badge variant={roleBadgeVariant(m.role)}>{roleLabel(m.role)}</Badge>
      ),
    },
    {
      key: "status",
      header: "Status",
      render: (m: WorkerCompany) => (
        <Badge variant={statusBadgeVariant(m.status)}>
          {m.status.charAt(0).toUpperCase() + m.status.slice(1)}
        </Badge>
      ),
    },
  ];

  return (
    <div>
      {/* Worker Info */}
      <PageHeader
        title={`${worker.firstName} ${worker.lastName}`}
        description={worker.email}
        action={
          <Button onClick={() => setEditModalOpen(true)}>Edit</Button>
        }
      />

      <div className="mb-8 rounded-lg border border-gray-200 bg-white p-6">
        <dl className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div>
            <dt className="text-sm font-medium text-gray-500">Email</dt>
            <dd className="mt-1 text-sm text-gray-900">{worker.email}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Telephone</dt>
            <dd className="mt-1 text-sm text-gray-900">{worker.phone ?? "\u2014"}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Auth Subject</dt>
            <dd className="mt-1 text-sm text-gray-900">{worker.authSubject}</dd>
          </div>
        </dl>
      </div>

      {/* Certificates Section */}
      <div className="mb-8">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900">Certificates</h2>
          <Button size="sm" onClick={() => setCertModalOpen(true)}>
            Add Certificate
          </Button>
        </div>
        <DataTable
          columns={certificateColumns}
          data={certificates}
          emptyMessage="No certificates found."
        />
      </div>

      {/* Company Memberships Section */}
      <div className="mb-8">
        <div className="mb-4">
          <h2 className="text-lg font-semibold text-gray-900">Company Memberships</h2>
        </div>
        <DataTable
          columns={membershipColumns}
          data={memberships}
          emptyMessage="No company memberships found."
        />
      </div>

      {/* Edit Worker Modal */}
      <Modal open={editModalOpen} onClose={() => setEditModalOpen(false)} title="Edit Worker">
        <form onSubmit={handleEditSubmit} className="space-y-4">
          <Input
            label="First Name"
            value={editForm.firstName}
            onChange={(e) => setEditForm({ ...editForm, firstName: e.target.value })}
            required
          />
          <Input
            label="Last Name"
            value={editForm.lastName}
            onChange={(e) => setEditForm({ ...editForm, lastName: e.target.value })}
            required
          />
          <Input
            label="Email"
            type="email"
            value={editForm.email}
            onChange={(e) => setEditForm({ ...editForm, email: e.target.value })}
            required
          />
          <Input
            label="Telephone"
            type="tel"
            value={editForm.phone}
            onChange={(e) => setEditForm({ ...editForm, phone: e.target.value })}
          />
          <Input
            label="Auth Subject"
            value={editForm.authSubject}
            onChange={(e) => setEditForm({ ...editForm, authSubject: e.target.value })}
            required
          />
          <div className="flex justify-end gap-3 pt-2">
            <Button type="button" variant="secondary" onClick={() => setEditModalOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={editSubmitting}>
              {editSubmitting ? "Saving..." : "Save Changes"}
            </Button>
          </div>
        </form>
      </Modal>

      {/* Add Certificate Modal */}
      <Modal open={certModalOpen} onClose={() => setCertModalOpen(false)} title="Add Certificate">
        <form onSubmit={handleCertSubmit} className="space-y-4">
          <Input
            label="Certificate Name"
            value={certForm.name}
            onChange={(e) => setCertForm({ ...certForm, name: e.target.value })}
            required
          />
          <Input
            label="Issuing Body"
            value={certForm.issuingBody}
            onChange={(e) => setCertForm({ ...certForm, issuingBody: e.target.value })}
          />
          <Input
            label="Certificate Number"
            value={certForm.certificateNumber}
            onChange={(e) => setCertForm({ ...certForm, certificateNumber: e.target.value })}
          />
          <Input
            label="Issued Date"
            type="date"
            value={certForm.issuedDate}
            onChange={(e) => setCertForm({ ...certForm, issuedDate: e.target.value })}
          />
          <Input
            label="Expiry Date"
            type="date"
            value={certForm.expiryDate}
            onChange={(e) => setCertForm({ ...certForm, expiryDate: e.target.value })}
          />
          <div className="flex justify-end gap-3 pt-2">
            <Button type="button" variant="secondary" onClick={() => setCertModalOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={certSubmitting}>
              {certSubmitting ? "Adding..." : "Add Certificate"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
