"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  PageHeader,
  DataTable,
  Badge,
  Button,
  Modal,
  Select,
} from "@/components/ui";
import { apiClient } from "@/lib/api/client";
import type {
  Shift,
  ShiftAssignment,
  Worksite,
  Worker,
} from "@/lib/types";

const statusVariant: Record<string, string> = {
  open: "info",
  assigned: "default",
  in_progress: "warning",
  completed: "success",
  cancelled: "danger",
};

const assignmentStatusVariant: Record<string, string> = {
  offered: "info",
  accepted: "success",
  declined: "danger",
  completed: "success",
};

function formatDate(iso: string | undefined): string {
  if (!iso) return "\u2014";
  const d = new Date(iso);
  return d.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export default function ShiftDetailPage() {
  const params = useParams();
  const router = useRouter();
  const shiftId = params.id as string;

  const [shift, setShift] = useState<Shift | null>(null);
  const [assignments, setAssignments] = useState<ShiftAssignment[]>([]);
  const [worksites, setWorksites] = useState<Worksite[]>([]);
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [assignModalOpen, setAssignModalOpen] = useState(false);
  const [selectedWorkerId, setSelectedWorkerId] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const worksiteLookup = worksites.reduce<Record<string, string>>(
    (acc, ws) => {
      acc[ws.id] = ws.name;
      return acc;
    },
    {}
  );

  const workerLookup = workers.reduce<Record<string, string>>(
    (acc, w) => {
      acc[w.id] = `${w.firstName} ${w.lastName}`;
      return acc;
    },
    {}
  );

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [shiftData, assignmentsData, worksitesData, workersData] =
        await Promise.all([
          apiClient<Shift>(`/api/v1/shifts/${shiftId}`),
          apiClient<ShiftAssignment[]>(`/api/v1/shifts/${shiftId}/assignments`),
          apiClient<Worksite[]>("/api/v1/worksites"),
          apiClient<Worker[]>("/api/v1/workers"),
        ]);
      setShift(shiftData);
      setAssignments(assignmentsData);
      setWorksites(worksitesData);
      setWorkers(workersData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  }, [shiftId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleStatusChange = async (newStatus: "completed" | "cancelled") => {
    setSubmitting(true);
    try {
      await apiClient(`/api/v1/shifts/${shiftId}/status`, {
        method: "PATCH",
        body: JSON.stringify({ status: newStatus }),
      });
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update status");
    } finally {
      setSubmitting(false);
    }
  };

  const handleAssignWorker = async () => {
    if (!selectedWorkerId) return;
    setSubmitting(true);
    try {
      await apiClient(`/api/v1/shifts/${shiftId}/assignments`, {
        method: "POST",
        body: JSON.stringify({ workerId: selectedWorkerId }),
      });
      setAssignModalOpen(false);
      setSelectedWorkerId("");
      await fetchData();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to assign worker"
      );
    } finally {
      setSubmitting(false);
    }
  };

  const handleAssignmentAction = async (
    assignmentId: string,
    action: "accept" | "decline"
  ) => {
    setSubmitting(true);
    try {
      await apiClient(
        `/api/v1/shifts/${shiftId}/assignments/${assignmentId}/${action}`,
        { method: "POST" }
      );
      await fetchData();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : `Failed to ${action} assignment`
      );
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="text-sm text-gray-500">Loading shift details...</div>;
  }

  if (error) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load shift: {error}
      </div>
    );
  }

  if (!shift) {
    return (
      <div className="text-sm text-gray-500">Shift not found.</div>
    );
  }

  const assignmentColumns = [
    {
      key: "workerId",
      header: "Worker",
      render: (a: ShiftAssignment) => workerLookup[a.workerId] || a.workerId,
    },
    {
      key: "status",
      header: "Status",
      render: (a: ShiftAssignment) => (
        <Badge
          variant={
            assignmentStatusVariant[a.status] as
              | "info"
              | "success"
              | "danger"
              | "default"
          }
        >
          {a.status}
        </Badge>
      ),
    },
    {
      key: "assignedAt",
      header: "Assigned At",
      render: (a: ShiftAssignment) => formatDate(a.assignedAt),
    },
    {
      key: "respondedAt",
      header: "Responded At",
      render: (a: ShiftAssignment) => formatDate(a.respondedAt),
    },
    {
      key: "actions",
      header: "Actions",
      render: (a: ShiftAssignment) =>
        a.status === "offered" ? (
          <div className="flex gap-2">
            <Button
              size="sm"
              variant="primary"
              disabled={submitting}
              onClick={(e) => {
                e.stopPropagation();
                handleAssignmentAction(a.id, "accept");
              }}
            >
              Accept
            </Button>
            <Button
              size="sm"
              variant="danger"
              disabled={submitting}
              onClick={(e) => {
                e.stopPropagation();
                handleAssignmentAction(a.id, "decline");
              }}
            >
              Decline
            </Button>
          </div>
        ) : null,
    },
  ];

  const workerOptions = workers.map((w) => ({
    value: w.id,
    label: `${w.firstName} ${w.lastName}`,
  }));

  return (
    <div>
      <div className="mb-6">
        <button
          onClick={() => router.push("/shifts")}
          className="text-sm text-blue-600 hover:text-blue-800"
        >
          &larr; Back to Shifts
        </button>
      </div>

      <PageHeader
        title={shift.title}
        description={shift.description || undefined}
        action={
          <div className="flex gap-2">
            {shift.status === "open" && (
              <Button
                variant="danger"
                disabled={submitting}
                onClick={() => handleStatusChange("cancelled")}
              >
                Cancel Shift
              </Button>
            )}
            {shift.status === "in_progress" && (
              <>
                <Button
                  variant="primary"
                  disabled={submitting}
                  onClick={() => handleStatusChange("completed")}
                >
                  Complete Shift
                </Button>
                <Button
                  variant="danger"
                  disabled={submitting}
                  onClick={() => handleStatusChange("cancelled")}
                >
                  Cancel Shift
                </Button>
              </>
            )}
          </div>
        }
      />

      <div className="mb-8 rounded-lg border border-gray-200 bg-white p-6">
        <h2 className="mb-4 text-lg font-semibold text-gray-900">
          Shift Details
        </h2>
        <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <dt className="text-sm font-medium text-gray-500">Worksite</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {worksiteLookup[shift.worksiteId] || "\u2014"}
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Status</dt>
            <dd className="mt-1">
              <Badge
                variant={
                  statusVariant[shift.status] as
                    | "info"
                    | "default"
                    | "warning"
                    | "success"
                    | "danger"
                }
              >
                {shift.status.replace("_", " ")}
              </Badge>
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Start Time</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {formatDate(shift.startTime)}
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">End Time</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {formatDate(shift.endTime)}
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Created At</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {formatDate(shift.createdAt)}
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Last Updated</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {formatDate(shift.updatedAt)}
            </dd>
          </div>
        </dl>
      </div>

      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-gray-900">Assignments</h2>
        <Button onClick={() => setAssignModalOpen(true)}>Assign Worker</Button>
      </div>
      <DataTable
        columns={assignmentColumns}
        data={assignments}
        emptyMessage="No workers assigned to this shift."
      />

      <Modal
        open={assignModalOpen}
        onClose={() => setAssignModalOpen(false)}
        title="Assign Worker"
      >
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleAssignWorker();
          }}
          className="space-y-4"
        >
          <Select
            label="Worker"
            options={workerOptions}
            value={selectedWorkerId}
            onChange={(e) => setSelectedWorkerId(e.target.value)}
          />
          <div className="flex justify-end gap-3 pt-2">
            <Button
              type="button"
              variant="secondary"
              onClick={() => setAssignModalOpen(false)}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={submitting || !selectedWorkerId}
            >
              {submitting ? "Assigning..." : "Assign Worker"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
