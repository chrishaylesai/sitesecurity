"use client";

import { useEffect, useState, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useAuth } from "@/lib/auth";

function CallbackHandler() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { handleCallback } = useAuth();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const process = async () => {
      const code = searchParams.get("code");
      const state = searchParams.get("state");
      const errorParam = searchParams.get("error");

      if (errorParam) {
        setError(
          `Authentication failed: ${searchParams.get("error_description") || errorParam}`
        );
        return;
      }

      if (!code || !state) {
        setError("Missing authorisation code or state parameter.");
        return;
      }

      try {
        await handleCallback(code, state);
        router.replace("/");
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Authentication failed"
        );
      }
    };

    process();
  }, [searchParams, handleCallback, router]);

  if (error) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="rounded-lg bg-red-50 p-6 text-center">
          <h2 className="text-lg font-semibold text-red-900">
            Authentication Error
          </h2>
          <p className="mt-2 text-sm text-red-700">{error}</p>
          <button
            onClick={() => router.push("/")}
            className="mt-4 rounded bg-red-600 px-4 py-2 text-sm text-white hover:bg-red-700"
          >
            Return to Dashboard
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <div className="text-center">
        <div className="mx-auto mb-4 h-10 w-10 animate-spin rounded-full border-4 border-gray-200 border-t-blue-600" />
        <p className="text-sm text-gray-600">Completing sign in...</p>
      </div>
    </div>
  );
}

export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-[60vh] items-center justify-center">
          <p className="text-sm text-gray-600">Loading...</p>
        </div>
      }
    >
      <CallbackHandler />
    </Suspense>
  );
}
