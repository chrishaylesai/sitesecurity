const variants: Record<string, string> = {
  default: "bg-gray-100 text-gray-700",
  success: "bg-gray-100 text-gray-700 ring-1 ring-inset ring-green-600/20",
  warning: "bg-gray-100 text-gray-700 ring-1 ring-inset ring-yellow-600/20",
  danger: "bg-gray-100 text-gray-700 ring-1 ring-inset ring-red-600/20",
  info: "bg-gray-100 text-gray-700 ring-1 ring-inset ring-blue-600/20",
};

interface BadgeProps {
  children: React.ReactNode;
  variant?: keyof typeof variants;
}

export function Badge({ children, variant = "default" }: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${variants[variant]}`}
    >
      {children}
    </span>
  );
}
